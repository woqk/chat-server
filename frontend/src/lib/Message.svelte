<script lang="ts">
  import { createEventDispatcher, onDestroy, onMount } from "svelte";
  import { fly } from "svelte/transition";
  import { format } from "date-fns";
  import Avatar from "./Avatar.svelte";
  import AppLocationRequest from "./AppLocationRequest.svelte";
  import { chat } from "../store";

  export let event: any;
  let replies: Array<Record<string, any>> = [];
  let isCurrentUserReplied: boolean = false;
  let isSenderCurrentUser: boolean = false;

  const dispatch = createEventDispatcher();

  function getComponent(name: string) {
    switch (name) {
      case "app.location.request":
        return AppLocationRequest;

      default:
        break;
    }
  }

  function onReply(evt: CustomEvent) {
    dispatch("reply", {
      ...evt.detail,
      replyTo: event.id,
    });
  }
  let unsubscribe;

  onMount(() => {
    // $0.scroll(0,$0.scrollHeight)
    dispatch("scrollBottom")
    const profile = chat.getProfile();
    isSenderCurrentUser = event.body.sender.id === profile.id;

    unsubscribe = chat.subscribe((value) => {
      replies = value.events.filter((evt) => {
        if (
          evt.type === "message" &&
          evt.body.sender.id === profile.id &&
          evt.body.replyTo === event.id
        ) {
          isCurrentUserReplied = true;
        }
        return evt.body.replyTo === event.id;
      });
    });
  });

  onDestroy(() => {
    unsubscribe && unsubscribe();
  });
</script>

<div transition:fly={{ y: 50, duration: 200 }} class="w-full flex my-4 text-sm">
  <div class="flex-shrink-0">
    <Avatar name={event.body.sender.username} />
  </div>
  <div class="ml-4">
    <div>
      <h5 class="text-white font-bold inline-block">
        {event.body.sender.username}
        {#if isSenderCurrentUser}
          <span class="text-gray-600 ml-1 text-xs"> (You) </span>
        {/if}
      </h5>
      <span class="mx-1 text-gray-500"> &middot; </span>
      <span class="text-gray-500 text-xs">
        {format(new Date(event.sentAt), "hh:mm bbb")}
      </span>
    </div>
    <div
      class="inline-block text-white break-all py-1.5 min-w-20 max-w-80% rounded-b-3xl rounded-l-3xl"
    >
      {#if event.body.type && event.body.type.startsWith("app.")}
        <svelte:component
          this={getComponent(event.body.type)}
          {event}
          {replies}
          {isCurrentUserReplied}
          on:reply={onReply}
        />
      {:else if event.type === "message"}
        {event.body.text}
      {:else if event.type === "login"}
        Successful login
      {/if}
    </div>
  </div>
</div>
