import { style } from "@vanilla-extract/css";
import Component from "/@/Components/Component";

export default () => {
  return (
    <div
      class={style({
        color: "red",
      })}
    >
      <Component />
    </div>
  );
};
