import {
  headerStyles,
  imgStyles,
  imgWrapperStyles,
  wrapperStyles,
} from "./Browser.css";

type BrowserProps = {
  url?: string;
  screenshotUrl: string;
};

export const Browser = ({
  url = "Not active",
  screenshotUrl,
}: BrowserProps) => {
  return (
    <div className={wrapperStyles}>
      <div className={headerStyles}>{url}</div>
      <div className={imgWrapperStyles}>
        <img src={screenshotUrl} className={imgStyles} />
      </div>
    </div>
  );
};

export default Browser;
