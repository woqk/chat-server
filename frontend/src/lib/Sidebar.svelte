<script lang="ts">
  import { active } from "tinro";
  import UserMenu from "./UserMenu.svelte";
  import {chat} from "../store"

  export let channels = [
    {
      name: "public",
    },
    {
      name: "go",
    },
    {
      name: "svelte",
    },
  ];
</script>

<div class="bg-gray-800 text-white p-2 overflow-auto flex flex-col">
  <h4 class="text-sm my-1 px-1 text-gray-200 font-bold">Channels</h4>
  <nav>
    {#each channels as channel}
      <a
        class="text-sm opacity-60 px-1 py-0.5 block hover:opacity-100"
        href="/channel/{channel.name}"
        use:active
      >
        # {channel.name}
      </a>
    {/each}
  </nav>
  <div class="block my-4" />
  <h4 class="text-sm my-1 px-1 text-gray-200 font-bold">Friends</h4>
  <nav>
    {#each channels as channel}
      <a
        class="text-sm opacity-60 px-1 py-0.5 block hover:opacity-100"
        href="/friend/12"
        use:active
      >
        # {channel.name}
      </a>
    {/each}
  </nav>
  <div class="mt-auto">

  </div>
  <UserMenu profile={$chat.profile} />
</div>
