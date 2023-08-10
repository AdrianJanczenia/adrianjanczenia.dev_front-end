import {useState} from "react";
import {useNavigate, useOutletContext} from "react-router-dom";
import "./css/Login.css";
import Input from "./form/Input";

const Login = () => {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    const {setJsonToken} = useOutletContext();

    const navigate = useNavigate();

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
                    // TODO wyswietlic error na stronie
                    console.log(data.message);
                } else {
                    setJsonToken(data.data.token);

                    // TODO wyswietlic alert na stronie

                    // toggleRefresh(true);
                    navigate("/home");
                }
            })
            .catch(error => {
                // TODO wyswietlic error na stronie
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
                        <p className="mt-5 mb-3 text-white text-center">&copy; Adrian Janczenia 2023</p>
                    </form>
                </main>
            </div>
        </>
    )
}

export default Login;