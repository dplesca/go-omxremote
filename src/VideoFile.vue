<template>
<div class="card video">
    <header class="card-header" @click="activecontrols = !activecontrols;">
        <p class="card-header-title">{{video.file}}</p>
    </header>
    <div class="card-content" v-show="activecontrols">
        <div class="content">
            <div class="tile">
                <a class="button is-primary is-outlined is-fullwidth" @click="handleClick('start', $event)"><span class="icon"><font-awesome-icon icon="play" /></span> <span>Start</span></a>
            </div>
            <div class="tile is-mobile">
                <div class="tile is-4">
                    <a class="button is-fullwidth" @click="handleClick('backward', $event)"><span class="icon is-small"><font-awesome-icon icon="backward" /></span> <span>Back</span></a>
                </div>
                <div class="tile is-4">
                    <a class="button is-fullwidth is-info is-outlined" @click="handleClick('pause', $event)"><span class="icon"><font-awesome-icon icon="pause" /></span> <span>Pause</span></a>
                </div>
                <div class="tile is-4">
                    <a class="button is-fullwidth" @click="handleClick('forward', $event)"><span class="icon is-small"><font-awesome-icon icon="forward" /></span> <span>Forward</span></a>
                </div>
            </div>
            <div class="tile"><a class="button is-fullwidth" @click="handleClick('subs', $event)"><span class="icon is-small"><font-awesome-icon icon="align-justify" /></span><span>Subs</span></a></div>
            <div class="tile"><a class="button is-primary is-danger is-fullwidth is-outlined" @click="handleClick('stop', $event)"><span class="icon"><font-awesome-icon icon="stop" /></span> <span>Stop</span></a></div>
        </div>
    </div>
</div>
</template>

<script>
import nanoajax from 'nanoajax'
import FontAwesomeIcon from '@fortawesome/vue-fontawesome'

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
    components: {
        FontAwesomeIcon
    }
}
</script>