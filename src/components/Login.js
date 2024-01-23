import {useEffect, useState} from "react";
import {useNavigate, useOutletContext} from "react-router-dom";
import "./css/Login.css";
import Input from "./form/Input";

const Login = () => {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    const {jsonToken} = useOutletContext();
    const {setJsonToken} = useOutletContext();
    const {Notification} = useOutletContext();

    const navigate = useNavigate();

    useEffect(() => {
        if (jsonToken !== "") {
            navigate("/home");
        }
    }, [jsonToken, navigate])

    const handleSubmit = (event) => {
        event.preventDefault();

        let payload = {
            action: "login",
            auth_payload: {
                username: username,
                password: password
            }
        }

        const requestOptions = {
            method: "POST",
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify(payload),
        }

        fetch(`${process.env.REACT_APP_BACKEND}/auth`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Notification.fire({
                        icon: 'error',
                        title: data.message,
                    });
                } else {
                    Notification.fire({
                        icon: 'success',
                        title: 'signed in successfully'
                    });

                    setTimeout(() => {
                        setJsonToken(data.data.token);
                        navigate("/home");
                    }, 2000);

                }
            })
            .catch(error => {
                console.log(error);
            })
    }

    return (
        <>
            <div className="container vertical-center">
                <main className="form-signin w-100 m-auto">
                    <form onSubmit={handleSubmit}>
                        <h1 className="h3 mb-3 fw-normal text-white">Sign in</h1>

                        <div className="form-floating">
                            <Input
                                type="text"
                                className="form-control"
                                id="username"
                                placeholder="Username"
                                onChange={(event) => setUsername(event.target.value)}
                                htmlFor="floatingInput"
                                titleLabel="Username"
                            />
                        </div>
                        <div className="form-floating">
                            <Input
                                type="password"
                                className="form-control"
                                id="password"
                                placeholder="Password"
                                onChange={(event) => setPassword(event.target.value)}
                                htmlFor="floatingInput"
                                titleLabel="Password"
                            />
                        </div>
                        <input className="btn btn-lg btn-light border-white w-100 py-2" type="submit" value="Sign in"/>
                        <p className="mt-5 mb-3 text-white text-center">&copy; Adrian Janczenia 2024</p>
                    </form>
                </main>
            </div>
        </>
    )
}

export default Login;