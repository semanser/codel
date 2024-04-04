import * as DropdownMenu from "@radix-ui/react-dropdown-menu";
import { useEffect } from "react";

import {
  Dropdown,
  dropdownMenuContentStyles,
  dropdownMenuItemStyles,
} from "@/components/Dropdown/Dropdown";
import { dropdownMenuItemIconStyles } from "@/components/Dropdown/Dropdown.css";
import { Icon } from "@/components/Icon/Icon";
import { Model } from "@/generated/graphql";

import { buttonStyles } from "./ModelSelector.css";

type ModelSelectorProps = {
  availableModels: Model[];
  selectedModel?: Model;
  activeModel?: Model;
  setSelectedModel: (model: Model) => void;
};

export const ModelSelector = ({
  availableModels = [],
  selectedModel,
  activeModel,
  setSelectedModel,
}: ModelSelectorProps) => {
  // Automatically select the first available model
  useEffect(() => {
    if (!activeModel && availableModels[0]) {
      setSelectedModel(availableModels[0]);
    }
  }, [activeModel, availableModels]);

  const handleValueChange = (value: string) => {
    const newModel = availableModels.find((model) => model.id === value);

    if (!newModel) return;

    setSelectedModel(newModel);
  };

  const dropdownContent = (
    <DropdownMenu.Content className={dropdownMenuContentStyles} sideOffset={5}>
      <DropdownMenu.RadioGroup
        value={selectedModel?.id}
        onValueChange={handleValueChange}
      >
        {availableModels.length > 0 ? (
          availableModels.map((model) => (
            <DropdownMenu.RadioItem
              key={model.id}
              className={dropdownMenuItemStyles}
              value={model.id}
            >
              <DropdownMenu.ItemIndicator
                className={dropdownMenuItemIconStyles}
              >
                <Icon.Check />
              </DropdownMenu.ItemIndicator>
              {model.id}
            </DropdownMenu.RadioItem>
          ))
        ) : (
          <DropdownMenu.Item disabled className={dropdownMenuItemStyles}>
            No available models
          </DropdownMenu.Item>
        )}
      </DropdownMenu.RadioGroup>
    </DropdownMenu.Content>
  );

  return (
    <Dropdown content={dropdownContent}>
      <div className={buttonStyles}>{activeModel?.id || "No model"}</div>
    </Dropdown>
  );
};
