<template>
<div class="tile is-ancestor is-vertical box video">
	<div class="tile">
		<h3 class="title is-4" @click="activecontrols = !activecontrols;">{{video.file}}</h3>
	</div>
	<div class="controls" v-show="activecontrols">
		<div class="tile">
			<a class="button is-primary is-outlined is-fullwidth" @click="handleClick('start', $event)"><span class="icon"><i class="fa fa-play-circle-o"></i></span> <span>Start</span></a>
		</div>
		<div class="tile is-mobile">
			<div class="tile is-4">
				<a class="button is-fullwidth" @click="handleClick('backward', $event)"><span class="icon is-small"><i class="fa fa-backward"></i></span> <span>Back</span></a>
			</div>
			<div class="tile is-4">
				<a class="button is-fullwidth is-info is-outlined" @click="handleClick('pause', $event)"><span class="icon"><i class="fa fa-pause-circle-o"></i></span> <span>Pause</span></a>
			</div>
			<div class="tile is-4">
				<a class="button is-fullwidth" @click="handleClick('forward', $event)"><span class="icon is-small"><i class="fa fa-forward"></i></span> <span>Forward</span></a>
			</div>
		</div>
		<div class="tile"><a class="button is-fullwidth" @click="handleClick('subs', $event)"><span class="icon is-small"><i class="fa fa-file-text-o"></i></span><span>Subs</span></a></div>
		<div class="tile"><a class="button is-primary is-danger is-fullwidth is-outlined" @click="handleClick('stop', $event)"><span class="icon"><i class="fa fa-stop-circle-o"></i></span> <span>Stop</span></a></div>
	</div>
</div>	
</template>

<script>
import nanoajax from 'nanoajax'
export default{
	name: "video-file",
	data(){
		return {
			activecontrols: false
		}
	},
	props: ['video'],
	methods: {
		handleClick(action, ev){
			ev.preventDefault();
			let requestURL = "/player/" + action;
			if (action == "start"){
				requestURL = "/start/" + this.video.hash;
			}
			nanoajax.ajax(
				{ url: requestURL, method: 'POST' },
				(code, responseText) => {})
		}
	},
}
</script>