<script>
import Icon from 'fa-svelte'
import { faPlay } from '@fortawesome/free-solid-svg-icons/faPlay'
import { faPause } from '@fortawesome/free-solid-svg-icons/faPause'
import { faStop } from '@fortawesome/free-solid-svg-icons/faStop'
import { faForward } from '@fortawesome/free-solid-svg-icons/faForward'
import { faBackward } from '@fortawesome/free-solid-svg-icons/faBackward'
import { faAlignJustify } from '@fortawesome/free-solid-svg-icons/faAlignJustify'
import Button from './Button.svelte'

export let file;
export let hash;

let activecontrols = false;

const clickVideo = (e) => {
    activecontrols = !activecontrols;
}

</script>

<div class="w-full rounded overflow-hidden shadow m-2">
    <header class="p-2" on:click={clickVideo}>
        <p class="font-semibold">{file}</p>
    </header>
    {#if activecontrols === true}
        <div class="flex m-1">
            <Button action="start" icon={faPlay} text="Start" hash="{hash}"></Button>
        </div>
        <div class="flex m-1">
            <Button action="backward" icon={faBackward} text="Back"></Button>
            <Button action="pause" icon={faPause} text="Pause"></Button>
            <Button action="forward" icon={faForward} text="Forward"></Button>
        </div>
        <div class="flex m-1">
            <Button action="subs" icon={faAlignJustify} text="Subs"></Button>
        </div>
        <div class="flex m-1">
            <Button action="stop" icon={faStop} text="Stop"></Button>
        </div>
    {/if}
</div>

<style>
header{cursor: pointer;}
</style>