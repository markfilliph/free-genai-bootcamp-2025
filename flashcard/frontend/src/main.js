import App from './App.svelte'
import { Router } from 'svelte-routing'

const app = new App({
    target: document.getElementById('app'),
    props: {
        Router
    }
})

export default app
