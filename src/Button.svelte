<script>
import Icon from 'fa-svelte'
import nanoajax from 'nanoajax';

export let icon;
export let action;
export let text;
export let hash = "";

let activestart = true;

const handleClick = (e) => {
  if (action == "start"){
    startVideo();
  } else {
    sendCommand(action);
  }
}

const startVideo = () => {
    if (activestart == true){
        activestart = false;
        nanoajax.ajax(
            { url: "/start/" + hash, method: 'POST' },
            (code, responseText) => {
                activestart = true;
        });
    }
}

const sendCommand = (action) => {
    let requestURL = "/player/" + action;
    nanoajax.ajax(
        { url: requestURL, method: 'POST' },
        (code, responseText) => {}
    );
}
</script>

<style type="text/postcss">
  .space {
    @apply px-4 py-2;
  }
  .outline {
    @apply bg-transparent font-semibold border rounded
  }
  .outline:hover{
    @apply text-white border-transparent
  }
  .pause {
    @apply text-blue-700 border-blue-500
  }
  .pause:hover{ 
    @apply bg-blue-500 
  }
  .start {
    @apply text-green-700 border-green-500
  }
  .start:hover{ 
    @apply bg-green-500 
  }
  .stop {
    @apply text-red-700 border-red-500
  }
  .stop:hover{ 
    @apply bg-red-500 
  }
  .backward, .forward, .prevsubs, .nextsubs {
    @apply text-gray-700 border-gray-500
  }
  .backward:hover, .forward:hover, .prevsubs:hover, .nextsubs:hover {
    @apply bg-gray-500
  }
</style>

<button
  class="space outline {action} w-full"
  on:click={handleClick}>
  <Icon icon={icon} /> <span>{text}</span>
</button>
