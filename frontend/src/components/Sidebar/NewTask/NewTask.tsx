import { useLocalStorage } from "@uidotdev/usehooks";
import { useNavigate } from "react-router-dom";

import { Tooltip } from "@/components/Tooltip/Tooltip";
import { Model } from "@/generated/graphql";

import { ModelSelector } from "./ModelSelector/ModelSelector";
import { linkWrapperStyles, wrapperStyles } from "./NewTask.css";

type NewTaskProps = {
  availableModels: Model[];
};

export const NewTask = ({ availableModels = [] }: NewTaskProps) => {
  const navigate = useNavigate();
  const [selectedModel, setSelectedModel] = useLocalStorage<Model | undefined>(
    "model",
  );
  const activeModel = availableModels.find(
    (model) => model.id == selectedModel?.id,
  );

  const handleNewTask = () => {
    navigate("/chat/new");
  };

  const tooltipContent = activeModel
    ? "Create a new flow"
    : "Please select a model first";

  return (
    <div className={wrapperStyles}>
      <Tooltip content={tooltipContent}>
        <button
          className={linkWrapperStyles}
          onClick={handleNewTask}
          disabled={!activeModel}
        >
          ✨ New task
        </button>
      </Tooltip>
      <ModelSelector
        availableModels={availableModels}
        selectedModel={selectedModel}
        activeModel={activeModel}
        setSelectedModel={setSelectedModel}
      />
    </div>
  );
};
