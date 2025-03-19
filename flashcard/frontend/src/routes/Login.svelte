<script>
    import { apiFetch } from '../lib/api.js'
    let email = ''
    let password = ''
    let error = null

    async function handleLogin() {
        try {
            const response = await apiFetch('/auth/login', {
                method: 'POST',
                body: JSON.stringify({ email, password })
            })
            // Handle successful login
        } catch (err) {
            error = err.message
        }
    }
</script>

<form on:submit|preventDefault={handleLogin}>
    <input type="email" bind:value={email} placeholder="Email" />
    <input type="password" bind:value={password} placeholder="Password" />
    <button type="submit">Login</button>
    {#if error}<div class="error">{error}</div>{/if}
</form>
