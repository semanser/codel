import { Link, useNavigate } from "react-router-dom";

import { wrapperStyles } from "./NewTask.css";

export const NewTask = () => {
  return (
    <Link to="/chat/new" className={wrapperStyles}>
      ✨ New task
    </Link>
  );
};
