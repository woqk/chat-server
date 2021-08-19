<script lang="ts">
  import clickOutside from "svelte-outside-click";
  import { createEventDispatcher } from "svelte";
  import Icon from "@iconify/svelte";
  import { fly } from "svelte/transition";

  export let channel: string = null;
  export let friend: string = null;
  let message: string;
  let appsOpen: boolean = false;
  let disableInput: boolean = false;

  const dispatch = createEventDispatcher();

  let apps = [
    {
      id: "app.location.request",
      icon: "lucide:map-pin",
      label: "Request Location",
      description: "Request location from users",
    },
    {
      id: "app.clipboard.request",
      icon: "lucide:clipboard",
      label: "Request Access Clipboard",
      description: "Access user clipboard",
    },
  ];

  function send() {
    dispatch("send", {
      type: "postMessage",
      body: {
        text: message,
        channel,
        user: friend,
      },
    });
    message = "";
  }

  function onKeyDown(evt: KeyboardEvent) {
    if (evt.key === "Enter") {
      send();
    }
  }

  function toggleApps() {
    appsOpen = !appsOpen;
  }

  function sendAppEvent(app: Record<string, string>) {
    disableInput = true;
    dispatch("send", {
      type: "postMessage",
      body: {
        type: app.id,
        text: message,
        channel,
        user: friend,
      },
    });
    appsOpen = false;
    disableInput = false;
    // setTimeout(() => {
    // }, 5000)
  }

  function onClickOutside() {
    appsOpen && (appsOpen = false);
  }
</script>

<div class="relative" use:clickOutside={onClickOutside}>
  {#if appsOpen}
    <div
      transition:fly={{ y: 50, duration: 200 }}
      class="bg-gray-700 h-250px absolute left-2 right-2 bottom-full rounded py-3 overflow-auto"
    >
      {#each apps as app}
        <div
          role="button"
          class="flex p-2 hover:(bg-gray-600)"
          aria-disabled={disableInput}
          on:click={() => !disableInput && sendAppEvent(app)}
        >
          <div
            class="bg-gray-800 text-white w-14 h-14 rounded-full flex justify-center items-center"
          >
            <Icon icon={app.icon} width={32} />
          </div>
          <div class="flex-grow ml-3 py-1.5">
            <h5 class="text-gray-200 text-base">{app.label}</h5>
            <p class="text-gray-400 text-sm">{app.description}</p>
          </div>
        </div>
      {/each}
    </div>
  {/if}

  <div
    class="h-12 flex justify-center items-center flex-shrink-0 m-2 bg-gray-800 text-sm rounded relative z-2 focus-within:bg-gray-700"
  >
    <button
      class="w-12 h-full flex justify-center items-center text-true-gray-300"
      on:click={() => !disableInput && toggleApps()}
      disabled={disableInput}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="24"
        height="24"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <line x1="12" y1="5" x2="12" y2="19" />
        <line x1="5" y1="12" x2="19" y2="12" />
      </svg>
    </button>
    <input
      class="flex-grow h-full flex-grow px-1 bg-transparent text-white outline-none"
      placeholder="Message here"
      on:keydown={onKeyDown}
      bind:value={message}
      disabled={disableInput}
    />
    <button
      class="w-12 pr-2 h-full flex justify-center items-center text-true-gray-300"
      on:click={() => !disableInput && send()}
      disabled={disableInput}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="24"
        height="24"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <line x1="22" y1="2" x2="11" y2="13" />
        <polygon points="22 2 15 22 11 13 2 9 22 2" />
      </svg>
    </button>
  </div>
</div>
