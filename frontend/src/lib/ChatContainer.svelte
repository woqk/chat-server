<script lang="ts">
  import { onMount } from "svelte";
  import Message from "./Message.svelte";
  import MessageCompose from "./MessageCompose.svelte";
  import { chat } from "../store";
  import WelcomeMessage from "./WelcomeMessage.svelte";

  export let channel: string = null;
  export let friend: string = null;
  let scrollContainer: HTMLElement;

  function send(evt: CustomEvent) {
    chat.send(evt.detail);
  }

  function sayHi() {
    chat.send({
      type: "postMessage",
      body: {
        text: "Hi 👋",
        channel,
        user: friend,
      },
    });
  }

  function onReply(evt: CustomEvent) {
    chat.send({
      type: "postMessage",
      body: {
        ...evt.detail,
        channel,
        user: friend,
      },
    });
  }

  onMount(() => {
    chat.send({
      type: "channel.join",
      body: {
        channel,
      },
    });
  });

  function scrollBottom() {
    scrollContainer && scrollContainer.scroll(0, scrollContainer.scrollHeight);
  }
</script>

<div class="w-full h-full overflow-hidden grid grid-rows-[1fr,auto]">
  <div
    class="overflow-auto flex items-end w-full h-full px-4 scrollbar-thin scrollbar-thumb-gray-500 scrollbar-track-gray-900 scrollbar-thumb-rounded-full"
    bind:this={scrollContainer}
    style="scroll-behavior: smooth"
  >
    <div class="mt-auto">
      <WelcomeMessage
        name={channel || friend}
        type={channel ? "channel" : "friend"}
        on:sayHi={sayHi}
      />
      {#each $chat.events as evt, index}
        {#if evt.type === "message" && evt.body.channel === channel && evt.body.replyTo === undefined}
          <Message
            event={evt}
            on:reply={onReply}
            on:scrollBottom={scrollBottom}
          />
        {/if}
      {/each}
    </div>
  </div>

  <MessageCompose {channel} {friend} on:send={send} />
</div>
