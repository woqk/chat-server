<script lang="ts">
  import clickOutside from "svelte-outside-click";
  import Avatar from "./Avatar.svelte";

  export let profile: Record<string, any>;

  let open: boolean = false;

  function onClickOutside() {
    open && (open = false);
  }
</script>

<details class="relative" bind:open use:clickOutside={onClickOutside}>
  <summary
    class="list-none w-full px-3 py-2 flex rounded-sm cursor-pointer select-none hover:bg-gray-700"
    class:bg-gray-700={open}
  >
    <Avatar name={profile.username} />
    <div class="ml-3 text-sm">
      {profile.username}
      <p class="text-gray-300 text-xs">{profile.id}</p>
    </div>
  </summary>
  <div
    class="absolute bottom-full mb-2 left-0 py-2 rounded-sm w-full bg-gray-700"
  >
    <nav>
      <a
        class="block py-2 px-1.5 hover:bg-gray-600"
        href="https://www.youtube.com/watch?v=dQw4w9WgXcQ"
        target="_blank"
        on:click={() => (open = false)}
      >
        Profile
      </a>
    </nav>
  </div>
</details>

<style>
  details > summary::-webkit-details-marker {
    display: none;
  }
</style>
