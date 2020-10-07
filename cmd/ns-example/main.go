package main

import (
	"context"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/go-git/go-git/v5"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/util/progress/progressui"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var (
	version  = "UNKNOWN"
	revision = "UNKNOWN"
)

func main() {
	log.Println(version, revision)
	example2()
	example()
}

func example() {
	/*
		Gitからビルド対象のリポジトリをクローン
	*/
	dir, err := ioutil.TempDir("", "clone")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	_, err = git.PlainClone(dir, false, &git.CloneOptions{URL: "https://github.com/yeasy/simple-web.git", Depth: 1})
	if err != nil {
		log.Fatal(err)
	}

	/*
		BuildkitdでリポジトリのDockerfileのビルドを実行し、生成されたDockerイメージをリポジトリにPush
	*/
	// buildkitdに接続
	c, err := client.New(context.Background(), "tcp://buildkitd:1234")
	if err != nil {
		log.Fatal(err)
	}

	ts := time.Now()
	dest := fmt.Sprintf("/data/image_build/%d", ts.Unix())
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// ビルドログの保存先作成
	logF, err := os.Create(filepath.Join(dest, "buildlog"))
	if err != nil {
		log.Fatal(err)
	}
	defer logF.Close()

	ch := make(chan *client.SolveStatus)
	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		// buildkitdでDockerfileを実行
		_, err := c.Solve(ctx, nil, client.SolveOpt{
			Exports: []client.ExportEntry{{
				Type: "image",
				Attrs: map[string]string{
					"name": "registry:5000/test/test",
					"push": "true",
				},
			}},
			LocalDirs: map[string]string{
				"context":    dir,
				"dockerfile": dir,
			},
			Frontend:      "dockerfile.v0",
			FrontendAttrs: map[string]string{"filename": "Dockerfile"},
		}, ch)
		return err
	})
	eg.Go(func() error {
		// ビルドログを収集
		return progressui.DisplaySolveStatus(context.Background(), "", nil, logF, ch)
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}

	/*
		ビルドしたイメージのコンテナを起動
	*/
	// Dockerデーモンに接続 (DooD)
	d, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	const imageName = "localhost:5000/test/test"

	// ビルドしたイメージをリポジトリからPull
	if err := d.PullImage(docker.PullImageOptions{
		Repository: imageName,
		Tag:        "latest",
	}, docker.AuthConfiguration{}); err != nil {
		log.Fatal(err)
	}

	// ビルドしたイメージのコンテナを作成
	container, err := d.CreateContainer(docker.CreateContainerOptions{
		Name: "ns_testcontainer", // コンテナ名
		Config: &docker.Config{
			Image: imageName,
			Labels: map[string]string{
				"neoshowcase.trap.jp/app":                                         "true",
				"traefik.enable":                                                  "true",
				"traefik.http.routers.ns_testcontainer.rule":                      "Host(`test.ns.localhost`)",
				"traefik.http.services.ns_testcontainer.loadbalancer.server.port": "80",
			},
		},
		HostConfig: &docker.HostConfig{
			RestartPolicy: docker.AlwaysRestart(),
		},
		NetworkingConfig: &docker.NetworkingConfig{EndpointsConfig: map[string]*docker.EndpointConfig{"neoshowcase_apps": {}}},
	})
	if err != nil {
		log.Fatal(err)
	}

	// コンテナを起動
	if err := d.StartContainer(container.ID, nil); err != nil {
		log.Fatal(err)
	}
}

func example2() {
	/*
		Gitからビルド対象のリポジトリをクローン
	*/
	dir, err := ioutil.TempDir("", "clone")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	_, err = git.PlainClone(dir, false, &git.CloneOptions{URL: "https://github.com/traPtitech/anke-to-UI.git", Depth: 1})
	if err != nil {
		log.Fatal(err)
	}

	/*
		ビルド用のBuildkitのLLBを構成
		このリポジトリはDockerfileが無い静的ファイルサイト
		ビルド後に必要な物のはdistディレクトリ内の静的ファイルのみ
	*/
	// ビルドステージを構成
	builder := llb.Image("docker.io/library/node:14.11.0-alpine"). // FROM node:14.11.0-alpine as builder
		Dir("/app"). // WORKDIR /app
		File(llb.Copy(llb.Local("local-src"), "package*.json", "./", &llb.CopyInfo{
			AllowWildcard:  true,
			CreateDestPath: true,
		})). // COPY package*.json ./
		Run(llb.Shlex("npm i")). // RUN npm i
		File(llb.Copy(llb.Local("local-src"), ".", ".", &llb.CopyInfo{
			AllowWildcard:  true,
			CreateDestPath: true,
		})). // COPY . .
		Run(llb.Shlex("npm run build"), llb.AddEnv("NODE_ENV", "production")). // RUN NODE_ENV=production npm run build
		Root()
	// ビルドで生成された静的ファイルのみを含むScratchイメージを構成
	def, err := llb.
		Scratch(). // FROM scratch
		File(llb.Copy(builder, "/app/dist", "/", &llb.CopyInfo{
			CopyDirContentsOnly: true,
			CreateDestPath:      true,
			AllowWildcard:       true,
		})). // COPY --from=builder /app/dist /
		Marshal(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	/*
		Buildkitdで構成したLLBを実行
	*/
	// buildkitdに接続
	c, err := client.New(context.Background(), "tcp://buildkitd:1234")
	if err != nil {
		log.Fatal(err)
	}

	ts := time.Now()
	dest := fmt.Sprintf("/data/static_build/%d", ts.Unix())
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// 成果物tarの保存先作成
	f, err := os.Create(filepath.Join(dest, "artifact.tar"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// ビルドログの保存先作成
	logF, err := os.Create(filepath.Join(dest, "buildlog"))
	if err != nil {
		log.Fatal(err)
	}
	defer logF.Close()

	ch := make(chan *client.SolveStatus)
	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		// buildkitdでLLBを実行
		_, err := c.Solve(ctx, def, client.SolveOpt{
			Exports: []client.ExportEntry{{
				Type:   "tar",
				Output: func(_ map[string]string) (io.WriteCloser, error) { return f, nil },
			}},
			LocalDirs: map[string]string{"local-src": dir},
		}, ch)
		return err
	})
	eg.Go(func() error {
		// ビルドログを収集
		return progressui.DisplaySolveStatus(context.Background(), "", nil, logF, ch)
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}
