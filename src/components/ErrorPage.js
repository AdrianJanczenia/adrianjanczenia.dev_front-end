import {useRouteError} from "react-router-dom";

export default function ErrorPage() {
    const error = useRouteError();

    return (
        <div className="container text-white">
            <div className="row vertical-center">
                <div className="col-md-3 offset-md-5">
                    <h1 className="mt-3">Oops!</h1>
                    <p>Unexpected error has occurred.</p>
                    <p>
                        <em>{error.statusText || error.message}</em>
                    </p>
                </div>
            </div>
        </div>
    )
}
