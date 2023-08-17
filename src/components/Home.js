import {useOutletContext} from "react-router-dom";

const Home = () => {
    const {jsonToken} = useOutletContext();
    const {setJsonToken} = useOutletContext();

    const logOut = () => {
        let payload = {
            action: "logout",
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
            .catch(error => {
                console.log("error logging out", error);
            })
            .finally(() => {
                setJsonToken("");
            })
    }

    return (
        <>
            <div className="text-center">
                <h2 style={{color: "white"}}>Home</h2>
            </div>
            <div className="text-center">
                {jsonToken !== "" &&
                    <a href="#!" onClick={logOut}><span className="badge bg-danger">Logout</span></a>
                }
            </div>
        </>
    )
}

export default Home;
