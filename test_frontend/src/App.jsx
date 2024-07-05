import './App.css'

function App() {
    const handleLogin = () => {
        window.location.href = "http://localhost:8080/auth/google/login";
    }

    return (
        <>
            <h1>Complete Auth</h1>

            <div className="card">
                <button onClick={handleLogin}>
                    Login with Google
                </button>
            </div>
        </>
    );
}

export default App;
