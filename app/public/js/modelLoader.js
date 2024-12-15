function loadModel(scene, loader, material, position, rotation, scale, modelPath, name, baseHeight) {
    return new Promise((resolve, reject) => {
        loader.load(modelPath, (geometry) => {
            const mesh = new THREE.Mesh(geometry, material);
            mesh.name = name;
            scene.add(mesh);
            console.log(name + ' model loaded');
            mesh.position.set(position.x, baseHeight + position.y, position.z);
            mesh.rotation.set(rotation.x, rotation.y, rotation.z);
            mesh.scale.set(scale.x, scale.y, scale.z);
            currentPositions[name.toLowerCase()] = { x: mesh.position.x, y: mesh.position.y, z: mesh.position.z };
            meshes[name.toLowerCase()] = mesh;
            const box = new THREE.Box3().setFromObject(mesh);
            const height = box.max.y - box.min.y;
            resolve(baseHeight + height);
        }, undefined, (error) => {
            console.error('An error occurred loading the ' + name + ' model:', error);
            reject(error);
        });
    });
}

async function loadModels(scene) {
    const loader = new THREE.STLLoader();
    
    const material = new THREE.MeshStandardMaterial({ color: 0xffffff, metalness: 0.3, roughness: 0.5 });

    const positions = [{ x: 15, y: -18, z: 0 }, { x: 15, y: -20, z: 0 }, { x: 15, y: -35, z: 0 }];
    const rotation = { x: 0, y: 0, z: Math.PI / 2 };
    const scale = { x: 0.1, y: 0.1, z: 0.1 };
    const models = ['base', 'fuselage', 'sensor_nose'];
    const names = ['Base', 'Fuselage', 'Nose'];

    // Load texture for base model
    const textureLoader = new THREE.TextureLoader();
    const baseTexture = textureLoader.load('public/textures/base_texture.png');
    const baseMaterial = new THREE.MeshStandardMaterial({ map: baseTexture, metalness: 0.3, roughness: 0.5 });

    try {
        let baseHeight = 0;
        for (let i = 0; i < models.length; i++) {
            const materialToUse = (names[i] === 'Base') ? baseMaterial : material;
            baseHeight = await loadModel(scene, loader, materialToUse, positions[i], rotation, scale, 'public/models/' + models[i] + '.stl', names[i], baseHeight);
        }
    } catch (error) {
        console.error('An error occurred while loading models:', error);
    }
}
