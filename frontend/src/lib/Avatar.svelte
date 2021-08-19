<script lang="ts">
  export let name: string;

  function getChars(): string {
    let chars: string = "";
    const f = name.split(" ");
    if (f.length >= 2) {
      chars += f[0][0];
      chars += f[1][0];
    } else {
      chars = f[0].substr(0, 2);
    }
    return chars;
  }

  $: chars = name ? getChars() : "no";
</script>

<div class="w-10 h-10 rounded-full bg-blue-300 uppercase flex justify-center items-center font-bold">
  {chars}
</div>
