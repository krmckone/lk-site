uniform vec2 u_resolution;
uniform vec2 u_mouse;
uniform float u_time;

void main() {
  vec2 st = gl_FragCoord.xy / u_resolution.xy;
  vec2 mouse = u_mouse.xy / u_resolution;
  gl_FragColor = vec4(st.x * mouse.x, st.y * mouse.y, 0.0, 1.0);
}