<script lang="ts">
  import { afterUpdate, onMount } from "svelte";
  import Avatar from "./Avatar.svelte";

  export let reply: Record<string, any>;
  let isFetched: boolean = false;
  let address: string = null;

  async function getAddress() {
    const data = await (
      await fetch(
        `https://nominatim.openstreetmap.org/reverse?lat=${reply.body.lat}&lon=${reply.body.lon}&format=json`
      )
    ).json();
    isFetched = true;
    address = data.display_name;
  }

  onMount(() => {});
  afterUpdate(async () => {
    if (reply && !reply.body.isDeny && address === null) {
      await getAddress();
    }
  });
  //   $: if (reply && !reply.body.isDeny && address === null) {
  //     getAddress();
  //   }
</script>

<li class="flex py-2">
  <div class="flex-shrink-0">
    <Avatar name={reply.body.sender.username} />
  </div>
  <div class="ml-3">
    <h5 class="font-bold">{reply.body.sender.username}</h5>
    {#if reply.body.isDeny}
      <p class="text-red-400">Denied</p>
    {:else if isFetched}
      <p class="text-gray-400">
        {address}
      </p>
    {:else}
      Fetching...
    {/if}
  </div>
</li>
