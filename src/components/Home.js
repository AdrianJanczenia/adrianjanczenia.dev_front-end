import {useEffect} from "react";
import {useNavigate, useOutletContext} from "react-router-dom";

const Home = () => {
    const {jsonToken} = useOutletContext();
    const navigate = useNavigate();

    useEffect(() => {
        console.log(jsonToken)
    }, [jsonToken, navigate])

    return (
        <>
            <div className="text-center">
                <h2>Home</h2>
            </div>
        </>
    )
}

export default Home;