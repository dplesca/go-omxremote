<script>
import Icon from 'fa-svelte'
import { faPlay } from '@fortawesome/free-solid-svg-icons/faPlay'
import { faPause } from '@fortawesome/free-solid-svg-icons/faPause'
import { faStop } from '@fortawesome/free-solid-svg-icons/faStop'
import Button from './Button.svelte'

export let link = "";

const onFocus = () => link = "";

</script>

<div class="w-full rounded overflow-hidden shadow p-2">
<header class="p-2">
    <input bind:value={link} on:focus={onFocus} type="text" placeholder="youtube link" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
</header>
<div class="flex my-1">
    <Button action="play" icon={faPlay} text="Play" link="{link}"></Button>
    <Button action="pause" icon={faPause} text="Pause"></Button>
    <Button action="stop" icon={faStop} text="Stop"></Button>
</div>
</div>