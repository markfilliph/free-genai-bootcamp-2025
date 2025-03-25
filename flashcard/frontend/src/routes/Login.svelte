<script>
    import { navigate, Link } from 'svelte-routing';
    import { apiFetch } from '../lib/api.js';
    
    let email = '';
    let password = '';
    let error = null;
    let isLoading = false;
    let emailError = '';
    let passwordError = '';
    
    // Simulate API fetch for demo purposes
    async function mockApiFetch(url, options) {
        return new Promise((resolve, reject) => {
            setTimeout(() => {
                // Demo credentials
                if (email === 'demo@example.com' && password === 'password123') {
                    resolve({
                        user: {
                            id: '123',
                            name: 'Demo User',
                            email: 'demo@example.com'
                        },
                        token: 'mock-jwt-token'
                    });
                } else {
                    reject(new Error('Invalid email or password'));
                }
            }, 1000); // Simulate network delay
        });
    }
    
    function validateForm() {
        let isValid = true;
        
        // Reset errors
        emailError = '';
        passwordError = '';
        error = null;
        
        // Validate email
        if (!email) {
            emailError = 'Email is required';
            isValid = false;
        } else if (!/^\S+@\S+\.\S+$/.test(email)) {
            emailError = 'Please enter a valid email address';
            isValid = false;
        }
        
        // Validate password
        if (!password) {
            passwordError = 'Password is required';
            isValid = false;
        } else if (password.length < 6) {
            passwordError = 'Password must be at least 6 characters';
            isValid = false;
        }
        
        return isValid;
    }
    
    async function handleLogin() {
        if (!validateForm()) return;
        
        isLoading = true;
        error = null;
        
        try {
            // Use mockApiFetch for demo, replace with real apiFetch in production
            const response = await mockApiFetch('/auth/login', {
                method: 'POST',
                body: JSON.stringify({ email, password })
            });
            
            // Store user data in localStorage
            localStorage.setItem('user', JSON.stringify(response.user));
            localStorage.setItem('token', response.token);
            
            // Navigate to home page
            navigate('/');
        } catch (err) {
            error = err.message;
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="login-container">
    <div class="login-card">
        <h2>Login to Your Account</h2>
        
        <form on:submit|preventDefault={handleLogin} aria-label="Login form" id="login-form">
            <div class="form-group">
                <label for="email">Email</label>
                <input 
                    type="email" 
                    id="email" 
                    bind:value={email} 
                    placeholder="Enter your email" 
                    aria-invalid={emailError ? 'true' : 'false'}
                    aria-describedby={emailError ? 'email-error' : undefined}
                />
                {#if emailError}
                    <div class="error email-error" id="email-error" aria-live="polite">{emailError}</div>
                {/if}
            </div>
            
            <div class="form-group">
                <label for="password">Password</label>
                <input 
                    type="password" 
                    id="password" 
                    bind:value={password} 
                    placeholder="Enter your password" 
                    aria-invalid={passwordError ? 'true' : 'false'}
                    aria-describedby={passwordError ? 'password-error' : undefined}
                />
                {#if passwordError}
                    <div class="error password-error" id="password-error" aria-live="polite">{passwordError}</div>
                {/if}
            </div>
            
            <button type="submit" class="login-button" disabled={isLoading}>
                {isLoading ? 'Logging in...' : 'Login'}
            </button>
            
            {#if error}
                <div class="error form-error" role="alert">{error}</div>
            {/if}
        </form>
        
        <div class="form-links">
            <Link to="/register">Create an account</Link>
            <Link to="/forgot-password">Forgot password?</Link>
        </div>
        
        <div class="demo-credentials">
            <p><strong>Demo Credentials:</strong></p>
            <p>Email: demo@example.com</p>
            <p>Password: password123</p>
        </div>
    </div>
</div>

<style>
    .login-container {
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 70vh;
    }
    
    .login-card {
        background: white;
        border-radius: 8px;
        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        padding: 2rem;
        width: 100%;
        max-width: 400px;
    }
    
    h2 {
        text-align: center;
        margin-bottom: 1.5rem;
        color: #333;
    }
    
    .form-group {
        margin-bottom: 1rem;
    }
    
    label {
        display: block;
        margin-bottom: 0.5rem;
        font-weight: 500;
    }
    
    input {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-size: 1rem;
    }
    
    input:focus {
        outline: none;
        border-color: #007bff;
        box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
    }
    
    input[aria-invalid="true"] {
        border-color: #dc3545;
    }
    
    .error {
        color: #dc3545;
        font-size: 0.875rem;
        margin-top: 0.25rem;
    }
    
    .login-button {
        width: 100%;
        padding: 0.75rem;
        background: #007bff;
        color: white;
        border: none;
        border-radius: 4px;
        font-size: 1rem;
        font-weight: 500;
        cursor: pointer;
        transition: background 0.3s ease;
    }
    
    .login-button:hover {
        background: #0069d9;
    }
    
    .login-button:disabled {
        background: #6c757d;
        cursor: not-allowed;
    }
    
    .form-links {
        display: flex;
        justify-content: space-between;
        margin-top: 1rem;
        font-size: 0.875rem;
    }
    
    .form-links :global(a) {
        color: #007bff;
        text-decoration: none;
    }
    
    .form-links :global(a:hover) {
        text-decoration: underline;
    }
    
    .demo-credentials {
        margin-top: 2rem;
        padding-top: 1rem;
        border-top: 1px solid #eee;
        font-size: 0.875rem;
    }
    
    .demo-credentials p {
        margin: 0.25rem 0;
    }
</style>
