import { MenuItem, MenuItemProps } from "./MenuItem/MenuItem";
import { NewTask } from "./NewTask/NewTask";
import { wrapperStyles } from "./Sidebar.css";

type SidebarProps = {
  items: MenuItemProps[];
  availableModels: string[];
};

export const Sidebar = ({ items = [], availableModels = [] }: SidebarProps) => {
  return (
    <div className={wrapperStyles}>
      <NewTask availableModels={availableModels} />
      {items.map((item) => (
        <MenuItem key={item.id} {...item} />
      ))}
    </div>
  );
};
