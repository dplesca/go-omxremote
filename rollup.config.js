import commonjs from 'rollup-plugin-commonjs';
// import purgeCss from '@fullhuman/postcss-purgecss';
import livereload from 'rollup-plugin-livereload';
import postcss from 'rollup-plugin-postcss';
import resolve from 'rollup-plugin-node-resolve';
import svelte from 'rollup-plugin-svelte';
import {terser} from 'rollup-plugin-terser';
import svelte_preprocess_postcss from 'svelte-preprocess-postcss';
import copy from 'rollup-plugin-copy';

const production = !process.env.ROLLUP_WATCH;
export default {
   input: 'src/main.js',
   output: {
      format: 'iife',
      sourcemap: false,
      name: 'app',
      file: 'dist/main.js',
   },
   plugins: [
      svelte({
         dev: !production,
         preprocess: {
            style: svelte_preprocess_postcss(),
         },
         css: css => {
            css.write('dist/components.css');
         },
      }),
      resolve(),
      commonjs(),
      postcss({
         sourceMap: false,
         extract: true,
      }),
      copy({
         targets: [
           { src: 'src/assets/*', dest: 'dist/assets' }
         ]
       }),
      !production && livereload('dist'),
      production && terser(),
   ],
   // watch: {
   //    clearScreen: false,
   // },
};
