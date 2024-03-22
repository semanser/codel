import { useEffect, useRef } from "react";

export const Browser = () => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const websocketRef = useRef<WebSocket>();

  useEffect(() => {
    if (!canvasRef.current || websocketRef.current) return;

    websocketRef.current = new WebSocket("ws://localhost:8080/stream");

    websocketRef.current.onmessage = (event) => {
      const canvas = canvasRef.current;
      if (!canvas) return;

      const context = canvas.getContext("2d");
      if (!context) return;

      const img = new Image();

      img.onload = function () {
        context.drawImage(img, 0, 0);
      };

      img.src = URL.createObjectURL(event.data);
    };
  }, []);

  return <canvas ref={canvasRef} width={640} height={480} />;
};

export default Browser;
