import {Outlet} from "react-router-dom";
import "./App.css";
import {useState} from "react";

function App() {
    const [jsonToken, setJsonToken] = useState("");

    return (
        <Outlet context={{jsonToken, setJsonToken}}/>
    );
}

export default App;
