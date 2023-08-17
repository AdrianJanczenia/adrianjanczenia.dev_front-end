import {Outlet, useNavigate} from "react-router-dom";
import "./App.css";
import {useEffect, useState} from "react";
import Swal from "sweetalert2";

function App() {
    const [jsonToken, setJsonToken] = useState("");
    const navigate = useNavigate();

    useEffect(() => {
        if (jsonToken === "") {
            let payload = {
                action: "refresh",
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
                        navigate("/");
                    } else {
                        setJsonToken(data.data.token);
                    }
                })
                .catch(error => {
                    console.log(error);
                })
        }
    }, [jsonToken, setJsonToken, navigate])

    const Notification = Swal.mixin({
        toast: true,
        position: 'top-end',
        showConfirmButton: false,
        timer: 2000,
        timerProgressBar: true,
        // willOpen: () => {
        //     setCSSLink();
        // },
        didOpen: (toast) => {
            toast.addEventListener('mouseenter', Swal.stopTimer)
            toast.addEventListener('mouseleave', Swal.resumeTimer)
        }
    })

    // TODO czy dark theme do alertow?
    // function setCSSLink() {
    //     let ss = document.createElement('link');
    //
    //     ss.rel = "stylesheet";
    //     ss.href = "https://cdn.jsdelivr.net/npm/@sweetalert2/theme-dark@5/dark.css";
    //
    //     document.head.appendChild(ss);
    // }

    return (
        <Outlet context={{jsonToken, setJsonToken, Notification}}/>
    );
}

export default App;
