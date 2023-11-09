import { styled } from "@macaron-css/solid";
import {
  FieldValue,
  FieldValues,
  FormStore,
  getValues,
} from "@modular-forms/solid";
import { relative } from "path";
import { Component, For, Show, createMemo, createSignal } from "solid-js";

const FieldValueList: Component<{
  values: FieldValues | (FieldValue | FieldValues)[];
  style?: string;
}> = (props) => {
  return (
    <ul style={props.style}>
      <For each={Object.entries(props.values)}>
        {([key, value]) => (
          <li>
            <span>{key}:</span>
            {!value ||
            typeof value !== "object" ||
            value instanceof File ||
            value instanceof Date ? (
              <span>{value instanceof File ? value.name : String(value)}</span>
            ) : (
              <FieldValueList style="padding-left: 1rem;" values={value} />
            )}
          </li>
        )}
      </For>
    </ul>
  );
};

const DebugFormStore: Component<{
  of: FormStore<any, any> | undefined;
}> = (props) => {
  const values = createMemo(() => props.of && getValues(props.of));

  return (
    <>
      state:
      <ul style="margin-bottom: 1rem;margin-left: 1rem;">
        <For
          each={[
            {
              label: "Submit Count",
              value: props.of?.submitCount,
            },
            {
              label: "Submitting",
              value: props.of?.submitting,
            },
            {
              label: "Submitted",
              value: props.of?.submitted,
            },
            {
              label: "Validating",
              value: props.of?.validating,
            },
            {
              label: "Dirty",
              value: props.of?.dirty,
            },
            {
              label: "Touched",
              value: props.of?.touched,
            },
            {
              label: "Invalid",
              value: props.of?.invalid,
            },
          ]}
        >
          {({ label, value }) => (
            <li>
              <span>{label}:</span>
              <span>{value?.toString()}</span>
            </li>
          )}
        </For>
      </ul>
      <Show
        when={Object.keys(values() || {}).length}
        fallback={<p>Wait for input...</p>}
      >
        values:
        <div style="margin-left:1rem;">
          <FieldValueList values={values()} />
        </div>
      </Show>
    </>
  );
};

export default DebugFormStore;
