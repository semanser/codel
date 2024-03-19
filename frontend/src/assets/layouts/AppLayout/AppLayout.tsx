import { Outlet } from "react-router-dom"
import { wrapperStyles } from "./AppLayout.css"
import { Sidebar } from "@/components/Sidebar/Sidebar"

const fakeData = [
  {
    title: "Test chat",
    id: 1
  },
  {
    title: "This is some very long title",
    id: 2
  },
  {
    title: "This is another chat",
    id: 3
  },
]

export const AppLayout = () => {
  return (
    <div className={wrapperStyles}>
      <Sidebar items={fakeData}/>
      <Outlet />
    </div>
  )
}
