import { Shape } from "./proto/generated/events";

export function drawRect(ctx: CanvasRenderingContext2D, shape: Shape) {
  ctx.fillStyle = shape.color;
  ctx.fillRect(shape.x, shape.y, shape.width, shape.height);
}

export function drawEllipse(ctx: CanvasRenderingContext2D, shape: Shape) {
  ctx.fillStyle = shape.color;
  ctx.beginPath();
  ctx.ellipse(
    shape.x + shape.width / 2,
    shape.y + shape.height / 2,
    shape.width / 2,
    shape.height / 2,
    0,
    0,
    Math.PI * 2,
  );
  ctx.fill();
}
