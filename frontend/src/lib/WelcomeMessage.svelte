<script lang="ts">
  import { createEventDispatcher } from "svelte";

  import { fly } from "svelte/transition";

  export let name: string = "";
  export let type: string = "";

  const dispatch = createEventDispatcher();

  function sayHi() {
    dispatch("sayHi", {});
  }
</script>

<div transition:fly={{ y: 50, duration: 200 }} class="w-full flex my-4 text-sm">
  <div class="ml-4">
    <div
      class="inline-block text-white break-all py-1.5 rounded-b-3xl rounded-l-3xl"
    >
      <h2 class="text-white font-bold text-3xl mb-3">
        {#if name && type === "channel"}
          #{name}
        {:else if name && type === "friend"}
          {name}
        {:else}
          Welcome
        {/if}
      </h2>
      {#if type === "friend"}
        <p class="text-gray-500">
          You'r in dm, this mean every things private.
        </p>
      {:else if type === "channel"}
        <button
          class="py-2 px-4 bg-gray-800 rounded-sm hover:opacity-80"
          on:click={sayHi}
        >
          Say Hi 👋
        </button>
      {/if}
      <!-- <p class="text-gray-500">Hi!</p> -->
    </div>
  </div>
</div>
