import {
  Navigate,
  Route,
  RouterProvider,
  createBrowserRouter,
  createRoutesFromElements,
} from "react-router-dom";

import { AppLayout } from "./assets/layouts/AppLayout/AppLayout";
import { ChatPage } from "./pages/ChatPage/ChatPage";
import "./styles/font.css.ts";
import "./styles/global.css.ts";
import "./styles/theme.css.ts";

export const router = createBrowserRouter(
  createRoutesFromElements(
    <>
      <Route element={<AppLayout />}>
        <Route path="/chat" element={<ChatPage />} />
      </Route>
      <Route path="*" element={<Navigate to="/chat" />} />
    </>,
  ),
);

function App() {
  return <RouterProvider router={router} />;
}

export default App;
