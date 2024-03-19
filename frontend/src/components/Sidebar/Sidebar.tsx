import { MenuItem, MenuItemProps } from "./MenuItem/MenuItem";
import { NewTask } from "./NewTask/NewTask";
import { wrapperStyles } from "./Sidebar.css";

type SidebarProps = {
  items: MenuItemProps[];
};

export const Sidebar = ({ items = [] }: SidebarProps) => {
  return (
    <div className={wrapperStyles}>
      <NewTask />
      {items.map((item) => (
        <MenuItem {...item} />
      ))}
    </div>
  );
};
