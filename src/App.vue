<template>
<div class="container">
	<div class="tile search">
		<div class="column">
			<div class="field">
				<p class="control has-icon has-icon-left is-expanded">
					<input type="text" class="input is-fullwidth" name="search" id="search" v-model="searchstring" @keyup="searchFiles"><span class="icon is-small is-left"><i class="fa fa-search"></i></span>
				</p>
			</div>
		</div>
	</div>
	<video-file v-for="item in files" :key="item.hash" :video="item" :showresult="item.show"></video-file>
</div>
</template>

<script>
import VideoFile from './VideoFile.vue'
import Fuse from 'fuse.js';
import nanoajax from 'nanoajax'

require("bulma/css/bulma.css")
require("font-awesome/css/font-awesome.css")
let fuse;
export default {
	name: 'app',
	data() {
		return {
			allFiles: [],
			files: [],
			searchstring: ""
		}
	},
	methods:{
		searchFiles(){
			if (this.searchstring){
				this.files = fuse.search(this.searchstring.toLowerCase());
			} else {
				this.files = this.allFiles;
			}
		}
	},
	mounted(){
		nanoajax.ajax(
			{ url:'/files.json' },
			(code, responseText) => {
				let files = JSON.parse(responseText);
				//files.forEach(function(element){ element.show = true;});
				this.files = files;
				this.allFiles = files;

				var options = {
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
				fuse = new Fuse(this.allFiles, options);
		})
	},
	components: {
		VideoFile
	}
}
</script>

<style>
body .title{
	font-weight: 400;
}
.tile, .video {
	margin: 5px 0;
}
.tile.is-ancestor:not(:last-child){
	margin-bottom: 1.25rem;
}
section.main{
	padding-top:1rem;
}
.search {
	margin:0 0 1.5rem 0;
}
header.card-header{
	cursor: pointer;
	min-width: 0;
}
.card-header p{
	overflow:hidden;
}
</style>
