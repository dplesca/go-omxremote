<script>
  import VideoFile from "./VideoFile.svelte";
  import Fuse from 'fuse.js';
	import { onMount } from "svelte";
  let videoFiles = [], files = [], q = "";
  let	fuse;

	onMount(async function() {
		const response = await fetch("/files.json");
		const json = await response.json();
    videoFiles = json;
    let options = {
      shouldSort: true,
      threshold: 0.6,
      location: 0,
      distance: 100,
      maxPatternLength: 32,
      minMatchCharLength: 1,
      keys: [
        "file"
      ]
    };
    fuse = new Fuse(videoFiles, options);
    if (q){
      files = fuse.search(q.toLowerCase())
    } else {
      files = videoFiles;
    }
  });
  
  const search = () => {
    if (q == ""){
      files = videoFiles;
    } else {
      files = fuse.search(q.toLowerCase());
    }
  }
</script>

<style>
header{background-image: linear-gradient(141deg,#009e6c,#00d1b2 71%,#00e7eb);}
</style>

<header class="px-20 py-5 w-full">
  <h1 class="text-3xl text-white">omxremote</h1>
</header>
<div class="container mx-auto my-5">
  <div class="p-2">
    <input bind:value={q} type="text" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" on:keyup={search}>
  </div>
	{#each files as videofile}
	<VideoFile file={videofile.file} hash={videofile.hash}></VideoFile>
	{/each}
</div>
