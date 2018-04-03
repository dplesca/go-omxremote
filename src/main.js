import Vue from 'vue'
import App from './App.vue'
import fontawesome from '@fortawesome/fontawesome'
import faPlay from '@fortawesome/fontawesome-free-solid/faPlay'
import faPause from '@fortawesome/fontawesome-free-solid/faPause'
import faForward from '@fortawesome/fontawesome-free-solid/faForward'
import faBackward from '@fortawesome/fontawesome-free-solid/faBackward'
import faStop from '@fortawesome/fontawesome-free-solid/faStop'
import faAlignJustify from '@fortawesome/fontawesome-free-solid/faAlignJustify'

fontawesome.library.add(faPlay, faPause, faForward, faBackward, faStop, faAlignJustify)

new Vue({
  el: '#app',
  render: h => h(App)
})
