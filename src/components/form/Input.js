import {forwardRef} from "react";

const Input = forwardRef((props, ref) => {
    return (
        <>
            <input
                type={props.type}
                className={props.className}
                id={props.id}
                placeholder={props.placeholder}
                onChange={props.onChange}
            />
            <label htmlFor={props.htmlFor} className="form-label">{props.titleLabel}</label>
        </>
    )
})

export default Input;
