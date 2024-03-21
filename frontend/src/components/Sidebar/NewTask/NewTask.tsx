import { NavLink } from "react-router-dom";

import { wrapperStyles } from "./NewTask.css";

export const NewTask = () => {
  return (
    <NavLink to="/chat/new" className={wrapperStyles}>
      ✨ New task
    </NavLink>
  );
};
