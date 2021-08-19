<script lang="ts">
  import Icon from "@iconify/svelte";
  import { createEventDispatcher } from "svelte";
  import AppLocationRequestLocationItem from "./AppLocationRequestLocationItem.svelte";
  import Avatar from "./Avatar.svelte";

  export let event: any;
  export let replies: Array<any> = [];
  let disableInput: boolean = false;
  let errorMessage: string = null;
  export let isCurrentUserReplied: boolean = false;

  const dispatch = createEventDispatcher();

  function success(position: GeolocationPosition): void {
    const latitude = position.coords.latitude;
    const longitude = position.coords.longitude;
    dispatch("reply", {
      isDeny: false,
      lat: latitude,
      lon: longitude,
    });
    disableInput = false;
  }

  function error(err: GeolocationPositionError): void {
    errorMessage = "Unable to retrieve your location";
    dispatch("reply", {
      isDeny: false,
      message: errorMessage,
    });
    disableInput = false;
  }

  function share() {
    disableInput = true;

    if (!navigator.geolocation) {
      errorMessage = "Geolocation is not supported by your browser";
      dispatch("reply", {
        isDeny: false,
        message: errorMessage,
      });
    } else {
      navigator.geolocation.getCurrentPosition(success, error);
    }
  }

  function deny() {
    dispatch("reply", {
      isDeny: true,
      message: "User deny",
    });
  }
  
</script>

<div class="rounded-sm bg-gray-700 w-300px py-6 px-5">
  <header class="text-center mb-5">
    <div class="mb-6 flex gap-1.5 justify-center items-center">
      <Icon icon="lucide:map-pin" width={36} />
      <span class="h-px w-4 border-b border-dashed border-gray-400 block" />
      <div
        class="rounded-full bg-gray-600 text-green-500 w-8 h-8 flex justify-center items-center"
      >
        <Icon icon="lucide:check" width={26} />
      </div>
      <span class="h-px w-4 border-b border-dashed border-gray-400 block" />
      <Avatar name={event.body.sender.username} />
    </div>
    <h3 class="text-lg mb-4">
      <span class="font-bold block">{event.body.sender.username}</span> wants your location.
    </h3>
    <!-- <p class="text-gray-400">Any description or note</p> -->
  </header>
  <section class="my-2">
    {#if replies.length > 0}
      <ul class="divide-y divide-gray-600 max-h-230px overflow-auto">
        {#each replies as reply}
          <AppLocationRequestLocationItem {reply} />
        {/each}
      </ul>
    {:else}
      <p class="text-center text-gray-400 py-4">Waiting</p>
    {/if}
  </section>
  {#if !isCurrentUserReplied}
    <footer class="flex gap-1.5">
      <button
        class="bg-red-500 py-2 px-4 rounded-sm flex-shrink"
        disabled={disableInput}
        on:click={deny}
      >
        Deny
      </button>
      <button
        class="bg-blue-500 py-2 px-4 rounded-sm flex-grow"
        disabled={disableInput}
        on:click={share}
      >
        Share
      </button>
    </footer>
  {/if}
  {#if errorMessage}
    <p class="text-red-400 my-3">{errorMessage}</p>
  {/if}
</div>
