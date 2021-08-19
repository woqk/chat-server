<script lang="ts">
  import Icon from "@iconify/svelte";

  import { onDestroy, onMount } from "svelte";
  import { Route } from "tinro";
  import {
    uniqueNamesGenerator,
    colors,
    animals,
  } from "unique-names-generator";
  import ChatContainer from "./lib/ChatContainer.svelte";
  import Sidebar from "./lib/Sidebar.svelte";
  import { chat } from "./store";

  let username: string;

  onMount(() => {
    username = uniqueNamesGenerator({
      dictionaries: [colors, animals],
      style: "lowerCase",
      separator: " ",
    });
    chat.connect();
  });

  $: if ($chat.status === "ready" && !$chat.isAuth) {
    chat.login(username);
  }

  onDestroy(() => {
    chat.disconnect();
  });
</script>

<main class="w-full h-screen bg-gray-900 grid grid-cols-[250px,1fr]">
  <Sidebar />
  {#if $chat.status === "ready" && $chat.isAuth}
    <Route path="/">
      <div class="w-full h-full flex flex-col items-center justify-center">
        <Icon class="text-white" icon="lucide:message-square" width={126} />
        <h2 class="text-3xl text-white mt-3">
          chat-server
        </h2>
      </div>
    </Route>
    <Route path="/channel/:channel" let:meta>
      <ChatContainer channel={meta.params.channel} />
    </Route>
    <Route path="/friend/:friend" let:meta>
      <ChatContainer friend={meta.params.friend} />
    </Route>
  {/if}
</main>

<!-- <main>
  <img src={logo} alt="Svelte Logo" />
  <h1>Hello Typescript!</h1>

  <Counter />

  <div class="flex flex-wrap gap-3 mx-5 my-10">
    <ChatContainer />
    <ChatContainer />
    <ChatContainer />
  </div>
</main>

<style>
  :root {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen,
      Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
  }
</style> -->
