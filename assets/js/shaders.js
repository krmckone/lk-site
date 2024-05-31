import * as THREE from 'three';
let container;
let camera, scene, renderer, clock;
let uniforms;

init();
animate();

async function init() {
  container = document.getElementById('container');

  camera = new THREE.Camera();
  camera.position.z = 1;

  scene = new THREE.Scene();
  clock = new THREE.Clock();

  const geometry = new THREE.PlaneGeometry(2, 2);

  uniforms = {
    u_time: { type: "f", value: 1.0 },
    u_resolution: { type: "v2", value: new THREE.Vector2() },
    u_mouse: { type: "v2", value: new THREE.Vector2() }
  };

  const vertexShader = await loadShaderSource('shader.vert')
  const fragmentShader = await loadShaderSource('shader.frag')

  const material = new THREE.ShaderMaterial(
    {
      uniforms,
      vertexShader,
      fragmentShader
    }
  );

  const mesh = new THREE.Mesh(geometry, material);
  scene.add(mesh);

  renderer = new THREE.WebGLRenderer();
  renderer.setPixelRatio(window.devicePixelRatio);

  container.appendChild(renderer.domElement);

  onWindowResize();
  window.addEventListener('resize', onWindowResize, false);

  document.onmousemove = function (e) {
    uniforms.u_mouse.value.x = e.pageX
    uniforms.u_mouse.value.y = e.pageY
  }
}

async function loadShaderSource(name) {
  const pathPrefix = "../shaders";
  const shaderResponse = await fetch(`${pathPrefix}/${name}`);
  const shaderSource = await shaderResponse.text();
  return shaderSource;
}

function onWindowResize(event) {
  renderer.setSize(container.clientWidth, container.clientHeight * 1.25);
  uniforms.u_resolution.value.x = renderer.domElement.width;
  uniforms.u_resolution.value.y = renderer.domElement.height;
}

function animate() {
  requestAnimationFrame(animate);
  render();
}

function render() {
  uniforms.u_time.value += clock.getDelta();
  if (renderer) {
    renderer.render(scene, camera);
  }
}