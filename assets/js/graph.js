

document.addEventListener('DOMContentLoaded', () => {
    // Basic setup from injected window variables
    const storyId = window.ForgeVerse?.storyId || 0;
    if (storyId === 0) {
        console.error("No story ID provided.");
        return;
    }

    // --- State ---
    const state = {
        zoom: 1,
        panX: 0,
        panY: 0,
        isPanning: false,
        startX: 0,
        startY: 0,
        nodes: [],       // Array of node objects: { id, name, type, x, y }
        edges: [],       // Array of edge objects: { id, from, to, text }
        isDraggingNode: false,
        draggedNode: null,
        connectionMode: false,
        connectionSourceNodeId: null,
        edgeBeingEdited: null,
        saveTimeout: null,
        _choiceRequired: false
    };

    // --- DOM Elements ---
    const graphContainer = document.getElementById('graphContainer');
    const graphGrid = document.getElementById('graphGrid');
    const graphNodesLayer = document.getElementById('graphNodesLayer');
    const graphSvg = document.getElementById('graphSvg');
    
    const zoomLevelDisplay = document.getElementById('zoomLevelDisplay');
    const saveStatus = document.getElementById('saveStatus');
    
    // --- Initialize ---
    function init() {
        setupEventListeners();
        loadGraph();
    }

    // --- Graph Operations (Load/Save) ---
    async function loadGraph() {
        try {
            const response = await fetch(`../api/load_graph.php?story_id=${storyId}`);
            if (!response.ok) throw new Error("Failed to load graph data");
            const data = await response.json();
            
            if (data.status === 'success') {
                state.nodes = data.nodes || [];
                state.edges = data.edges || [];
                
                // If completely new story, initialize a Start Node
                if (state.nodes.length === 0) {
                    createNode('Start', 'start', 100, window.innerHeight / 2 - 100);
                    // auto save will catch it
                }
                
                document.getElementById('storyNameDisplay').innerText = data.story_name || "Story Graph";
                
                renderAll();
                centerGraph();
            } else {
                console.error(data.message);
                alert("Error loading graph: " + data.message);
            }
        } catch (err) {
            console.error(err);
        }
    }

    function autoSave() {
        if (state.saveTimeout) clearTimeout(state.saveTimeout);
        
        saveStatus.classList.remove('visible');
        saveStatus.innerText = "Saving...";
        saveStatus.style.color = "var(--text-secondary)";
        saveStatus.classList.add('visible');

        state.saveTimeout = setTimeout(async () => {
            try {
                const response = await fetch('../api/save_graph.php', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        story_id: storyId,
                        nodes: state.nodes,
                        edges: state.edges
                    })
                });
                
                const data = await response.json();
                if (data.status === 'success') {
                    // update ids if new
                    if (data.idMap) {
                        // Backend can return mapping for newly inserted items if needed
                        // (Usually simplified by using unique client-side IDs for new items before sync)
                    }
                    saveStatus.innerText = "Saved ✔";
                    saveStatus.style.color = "var(--accent-green)";
                    setTimeout(() => saveStatus.classList.remove('visible'), 2000);
                }
            } catch (err) {
                console.error("AutoSave Error:", err);
                saveStatus.innerText = "Save Failed!";
                saveStatus.style.color = "var(--accent-red)";
            }
        }, 1000); // Debounce save
    }

    // --- Core Rendering ---
    function renderAll() {
        updateTransform();
        renderNodes();
        renderEdges();
    }

    function renderNodes() {
        graphNodesLayer.innerHTML = '';
        state.nodes.forEach(node => {
            const el = document.createElement('div');
            el.className = `graph-node node-${node.type}`;
            el.id = `node-${node.id}`;
            el.style.left = `${node.x}px`;
            el.style.top = `${node.y}px`;
            
            // Calculate status
            let statusHtml = '';
            if (node.type === 'scene') {
                const incomingEdges = state.edges.filter(e => e.to == node.id);
                if (incomingEdges.length === 0) {
                    statusHtml = `<div class="node-status status-red" title="Unreachable. Click to connect." data-action="fix-unreachable">✕</div>`;
                } else {
                    statusHtml = `<div class="node-status status-green" title="Connected">✔</div>`;
                }
            }

            let icon = '🎬';
            if (node.type === 'start') icon = '🟢';
            if (node.type === 'end') icon = '🛑';

            let bodyHtml = '';
            if (node.type === 'start') {
                const hasOutgoing = state.edges.some(e => e.from == node.id);
                const startStatus = !hasOutgoing
                    ? `<div class="node-status status-red" title="Start is not connected to any scene. Click to fix." data-action="fix-start">✕</div>`
                    : `<div class="node-status status-green" title="Connected">✔</div>`;
                bodyHtml = `
                <div class="node-body">
                    <div class="node-actions">
                        <button class="node-btn add-btn" title="Create Connected Scene" data-action="add">+</button>
                        <button class="node-btn link-btn" title="Connect to Existing Scene" data-action="link">🔗</button>
                    </div>
                    ${startStatus}
                </div>`;
            } else if (node.type === 'scene') {
                bodyHtml = `
                <div class="node-body">
                    <div class="node-actions">
                        <button class="node-btn add-btn" title="Create Connected Scene" data-action="add">+</button>
                        <button class="node-btn delete-btn" title="Delete Scene" data-action="delete">🗑</button>
                        <button class="node-btn link-btn" title="Connect Scene" data-action="link">🔗</button>
                    </div>
                    ${statusHtml}
                </div>`;
            }

            el.innerHTML = `
                <div class="node-header">
                    <span class="node-icon">${icon}</span>
                    <span class="node-title" title="Open Scene Editor">${node.name}</span>
                </div>
                ${bodyHtml}
                ${node.type !== 'start' ? `<div class="node-port port-in"></div>` : ''}
                ${node.type !== 'end' ? `<div class="node-port port-out"></div>` : ''}
            `;
            
            // Event listeners for node
            el.addEventListener('mousedown', (e) => handleNodeMouseDown(e, node));
            
            // Action listeners
            const title = el.querySelector('.node-title');
            if(title) {
                if (node.type === 'start') {
                    // Start node is void — no editor, non-interactive title
                    title.classList.add('node-title-void');
                } else {
                    title.addEventListener('click', (e) => {
                        e.stopPropagation();
                        window.location.href = `editor.php?story_id=${storyId}&scene_id=${node.id}`;
                    });
                }
            }

            const addBtn = el.querySelector('[data-action="add"]');
            if(addBtn) {
                addBtn.addEventListener('click', (e) => {
                    e.stopPropagation();
                    state.parentNodeIdForNewScene = node.id;

                    // Show/hide choice text field based on node type & existing connections
                    const existingOut = state.edges.filter(ed => ed.from == node.id).length;
                    const isStart = node.type === 'start';
                    const choiceGroup   = document.getElementById('newSceneChoiceGroup');
                    const choiceReqSpan = document.getElementById('newSceneChoiceRequired');
                    const choiceOptSpan = document.getElementById('newSceneChoiceOptional');

                    if (isStart) {
                        choiceGroup.style.display = 'none';
                    } else if (existingOut === 0) {
                        choiceGroup.style.display   = 'block';
                        choiceReqSpan.style.display = 'none';
                        choiceOptSpan.style.display = 'inline';
                    } else {
                        choiceGroup.style.display   = 'block';
                        choiceReqSpan.style.display = 'inline';
                        choiceOptSpan.style.display = 'none';
                    }

                    openModal('addScene');
                });
            }

            const delBtn = el.querySelector('[data-action="delete"]');
            if(delBtn) {
                delBtn.addEventListener('click', (e) => {
                    e.stopPropagation();
                    deleteNode(node.id);
                });
            }

            const linkBtn = el.querySelector('[data-action="link"]');
            if(linkBtn) {
                linkBtn.addEventListener('click', (e) => {
                    e.stopPropagation();
                    startConnectionMode(node.id);
                });
            }

            const fixBtn = el.querySelector('[data-action="fix-unreachable"]');
            if(fixBtn) {
                fixBtn.addEventListener('click', (e) => {
                    e.stopPropagation();
                    alert("This scene is unreachable. Use the 🔗 Link button on another scene to connect to this one.");
                });
            }

            const fixStartBtn = el.querySelector('[data-action="fix-start"]');
            if(fixStartBtn) {
                fixStartBtn.addEventListener('click', (e) => {
                    e.stopPropagation();
                    startConnectionMode(node.id);
                });
            }

            graphNodesLayer.appendChild(el);
        });
    }

    function renderEdges() {
        // Clear old edges (keep defs)
        const defs = graphSvg.querySelector('defs').outerHTML;
        graphSvg.innerHTML = defs;

        state.edges.forEach(edge => {
            const fromNode = state.nodes.find(n => n.id == edge.from);
            const toNode = state.nodes.find(n => n.id == edge.to);
            if (!fromNode || !toNode) return;

            // Get port positions (relative to container)
            // Node width 200, header ~45. Port Out is bottom center. Port In is top center.
            const fromX = fromNode.x + 100;
            const fromY = fromNode.y + (fromNode.type === 'scene' ? 105 : 45); // Approximate height
            
            const toX = toNode.x + 100;
            const toY = toNode.y;

            // Create SVG path (bezier curve)
            const d = getBezierPath(fromX, fromY, toX, toY);
            
            // Hitbox
            const hitbox = document.createElementNS("http://www.w3.org/2000/svg", "path");
            hitbox.setAttribute('d', d);
            hitbox.setAttribute('class', 'connection-hitbox');
            
            // Visible line
            const path = document.createElementNS("http://www.w3.org/2000/svg", "path");
            path.setAttribute('d', d);
            path.setAttribute('class', 'connection-line');
            path.setAttribute('marker-end', 'url(#arrowhead)');
            
            // Click to edit
            const clickHandler = (e) => {
                e.stopPropagation();
                openEditConnectionModal(edge);
            };
            hitbox.addEventListener('click', clickHandler);
            path.addEventListener('click', clickHandler);

            // Hover effects
            const mouseEnter = () => { path.style.stroke = 'var(--accent-blue)'; path.setAttribute('marker-end', 'url(#arrowhead-highlight)'); };
            const mouseLeave = () => { path.style.stroke = ''; path.setAttribute('marker-end', 'url(#arrowhead)'); };
            hitbox.addEventListener('mouseenter', mouseEnter);
            hitbox.addEventListener('mouseleave', mouseLeave);

            graphSvg.appendChild(hitbox);
            graphSvg.appendChild(path);

            // Label — only show for non-start connections that have choice text
            const edgeText = (edge.text || '').trim();
            if (fromNode.type !== 'start' && edgeText) {
                const midX = (fromX + toX) / 2;
                const midY = (fromY + toY) / 2;

                const label = document.createElement('div');
                label.className = 'connection-label';
                label.innerText = edgeText;
                label.style.left = `${midX}px`;
                label.style.top = `${midY}px`;
                label.addEventListener('click', clickHandler);
                label.addEventListener('mouseenter', mouseEnter);
                label.addEventListener('mouseleave', mouseLeave);

                graphNodesLayer.appendChild(label);
            }
        });

        // If drawing line active
        if (state.connectionMode && state.drawingLine) {
            graphSvg.appendChild(state.drawingLine);
        }
    }

    function getBezierPath(x1, y1, x2, y2) {
        const dy = Math.abs(y2 - y1);
        const curveOffset = Math.max(50, dy * 0.5);
        return `M ${x1} ${y1} C ${x1} ${y1 + curveOffset}, ${x2} ${y2 - curveOffset}, ${x2} ${y2}`;
    }

    // --- Interactions: Pan & Zoom ---
    function updateTransform() {
        const transform = `translate(${state.panX}px, ${state.panY}px) scale(${state.zoom})`;
        graphNodesLayer.style.transform = transform;
        graphSvg.style.transform = transform;
        
        // Background Grid (simulate infinite by background-position)
        graphGrid.style.transform = `scale(${state.zoom})`;
        graphGrid.style.backgroundPosition = `${state.panX / state.zoom}px ${state.panY / state.zoom}px`;
        
        zoomLevelDisplay.innerText = `${Math.round(state.zoom * 100)}%`;
    }

    function centerGraph() {
        if (state.nodes.length === 0) return;
        let minX = Infinity, minY = Infinity, maxX = -Infinity, maxY = -Infinity;
        state.nodes.forEach(n => {
            minX = Math.min(minX, n.x);
            minY = Math.min(minY, n.y);
            maxX = Math.max(maxX, n.x + 200);
            maxY = Math.max(maxY, n.y + 100);
        });
        
        const centerX = (minX + maxX) / 2;
        const centerY = (minY + maxY) / 2;
        
        const containerRect = graphContainer.getBoundingClientRect();
        
        state.zoom = 1;
        state.panX = (containerRect.width / 2) - centerX;
        state.panY = (containerRect.height / 2) - centerY;
        
        updateTransform();
    }

    // --- Interactions: Drag Node ---
    function handleNodeMouseDown(e, node) {
        if (state.connectionMode) {
            // End connection mode by selecting target
            completeConnection(node.id);
            return;
        }

        if (e.target.closest('.node-actions') || e.target.closest('.node-title') || e.target.closest('.node-status')) {
            return; // Don't drag if clicking buttons
        }

        e.stopPropagation();
        state.isDraggingNode = true;
        state.draggedNode = node;
        
        // Convert screen coordinates to canvas coordinates considering zoom
        state.startX = e.clientX / state.zoom - node.x;
        state.startY = e.clientY / state.zoom - node.y;

        document.getElementById(`node-${node.id}`).classList.add('selected');
    }

    // --- Event Listeners Setup ---
    function setupEventListeners() {
        // Container Pan / Drag
        graphContainer.addEventListener('mousedown', (e) => {
            if (e.button !== 0) return; // Only left click
            if (state.connectionMode) return;
            
            state.isPanning = true;
            state.startX = e.clientX - state.panX;
            state.startY = e.clientY - state.panY;
        });

        window.addEventListener('mousemove', (e) => {
            if (state.isDraggingNode && state.draggedNode) {
                state.draggedNode.x = e.clientX / state.zoom - state.startX;
                state.draggedNode.y = e.clientY / state.zoom - state.startY;
                
                // Snap to grid (optional, snap 20px)
                state.draggedNode.x = Math.round(state.draggedNode.x / 20) * 20;
                state.draggedNode.y = Math.round(state.draggedNode.y / 20) * 20;

                renderAll(); // Rerender to update lines immediately
            } 
            else if (state.isPanning) {
                state.panX = e.clientX - state.startX;
                state.panY = e.clientY - state.startY;
                updateTransform();
            }
            else if (state.connectionMode && state.drawingLine) {
                // Update drawing line
                const fromNode = state.nodes.find(n => n.id == state.connectionSourceNodeId);
                if(fromNode) {
                    const fromX = fromNode.x + 100;
                    const fromY = fromNode.y + (fromNode.type === 'scene' ? 105 : 45);
                    
                    const containerRect = graphContainer.getBoundingClientRect();
                    const mouseX = (e.clientX - containerRect.left - state.panX) / state.zoom;
                    const mouseY = (e.clientY - containerRect.top - state.panY) / state.zoom;
                    
                    state.drawingLine.setAttribute('d', getBezierPath(fromX, fromY, mouseX, mouseY));
                }
            }
        });

        window.addEventListener('mouseup', () => {
            if (state.isDraggingNode) {
                state.isDraggingNode = false;
                if(state.draggedNode) {
                    document.getElementById(`node-${state.draggedNode.id}`)?.classList.remove('selected');
                    state.draggedNode = null;
                    autoSave();
                }
            }
            state.isPanning = false;
        });

        // Zoom
        graphContainer.addEventListener('wheel', (e) => {
            e.preventDefault();
            const zoomFactor = 0.1;
            const delta = e.deltaY > 0 ? -1 : 1;
            
            // Calculate mouse position relative to container
            const containerRect = graphContainer.getBoundingClientRect();
            const mouseX = e.clientX - containerRect.left;
            const mouseY = e.clientY - containerRect.top;

            // Calculate point in canvas coordinates before zoom
            const pointX = (mouseX - state.panX) / state.zoom;
            const pointY = (mouseY - state.panY) / state.zoom;

            // Update zoom
            state.zoom += delta * zoomFactor;
            state.zoom = Math.max(0.2, Math.min(state.zoom, 3)); // clamp 20% to 300%

            // Adjust pan to zoom towards mouse cursor
            state.panX = mouseX - pointX * state.zoom;
            state.panY = mouseY - pointY * state.zoom;

            updateTransform();
        }, { passive: false });

        // Toolbar Buttons
        document.getElementById('zoomInBtn').addEventListener('click', () => { state.zoom = Math.min(3, state.zoom + 0.2); updateTransform(); });
        document.getElementById('zoomOutBtn').addEventListener('click', () => { state.zoom = Math.max(0.2, state.zoom - 0.2); updateTransform(); });
        document.getElementById('centerGraphBtn').addEventListener('click', centerGraph);
        document.getElementById('backBtn').addEventListener('click', () => window.location.href = 'dashboard.php');

        // Modals Setup
        const resetAddModal = () => {
            closeModal('addScene');
            state.parentNodeIdForNewScene = null;
            const ci = document.getElementById('newSceneChoice');
            if (ci) { ci.value = ''; ci.style.borderColor = ''; ci.placeholder = 'What choice leads to this scene?'; }
        };
        document.getElementById('addSceneClose').addEventListener('click', resetAddModal);
        document.getElementById('addSceneBackdrop').addEventListener('click', resetAddModal);
        
        document.getElementById('createSceneBtn').addEventListener('click', () => {
            const name = document.getElementById('newSceneName').value.trim();
            if (!name) return;

            const choiceInput  = document.getElementById('newSceneChoice');
            const choiceText   = choiceInput ? choiceInput.value.trim() : '';
            const parentNode   = state.nodes.find(n => n.id == state.parentNodeIdForNewScene);

            // Validate choice text if required
            if (parentNode) {
                const isStart        = parentNode.type === 'start';
                const existingOut    = state.edges.filter(e => e.from == parentNode.id).length;
                const choiceRequired = !isStart && existingOut >= 1;

                if (choiceRequired && !choiceText) {
                    choiceInput.style.borderColor = 'var(--accent-red)';
                    choiceInput.placeholder = 'Choice text is required when multiple connections exist';
                    return;
                }
                if (choiceInput) choiceInput.style.borderColor = '';
            }

            let x, y;
            if (parentNode) {
                x = parentNode.x + 250;
                y = parentNode.y;
                while (state.nodes.some(n => Math.abs(n.x - x) < 50 && Math.abs(n.y - y) < 50)) {
                    y += 120;
                }
            } else {
                const rect = graphContainer.getBoundingClientRect();
                x = ((rect.width / 2) - state.panX) / state.zoom - 100;
                y = ((rect.height / 2) - state.panY) / state.zoom - 50;
            }

            const newNodeId = generateId();
            const newNode = {
                id: newNodeId,
                name,
                type: 'scene',
                x: Math.round(x / 20) * 20,
                y: Math.round(y / 20) * 20
            };
            state.nodes.push(newNode);

            if (parentNode) {
                const isStart   = parentNode.type === 'start';
                const edgeLabel = isStart ? '' : choiceText;  // Start node: never a choice label
                state.edges.push({
                    id:   generateId(),
                    from: parentNode.id,
                    to:   newNodeId,
                    text: edgeLabel
                });
            }

            state.parentNodeIdForNewScene = null;
            renderAll();
            autoSave();

            closeModal('addScene');
            document.getElementById('newSceneName').value = '';
            if (choiceInput) { choiceInput.value = ''; choiceInput.style.borderColor = ''; choiceInput.placeholder = 'What choice leads to this scene?'; }
        });

        // Connect Scene Modal
        const resetConnectModal = () => {
            closeModal('connectScene');
            cancelConnectionMode();
            const ci = document.getElementById('connectionText');
            if (ci) { ci.style.borderColor = ''; ci.placeholder = 'What choice leads here?'; }
            state._choiceRequired = false;
        };
        document.getElementById('connectSceneClose').addEventListener('click', resetConnectModal);
        document.getElementById('connectSceneBackdrop').addEventListener('click', resetConnectModal);
        document.getElementById('confirmConnectBtn').addEventListener('click', () => {
            const connectionInput = document.getElementById('connectionText');
            const text = connectionInput.value.trim();

            if (state._choiceRequired && !text) {
                connectionInput.style.borderColor = 'var(--accent-red)';
                connectionInput.placeholder = 'Choice text is required when multiple connections exist';
                return;
            }
            connectionInput.style.borderColor = '';

            createEdge(state.connectionSourceNodeId, state.connectionTargetNodeId, text);
            closeModal('connectScene');
            state.connectionSourceNodeId = null;
            state.connectionTargetNodeId = null;
            state._choiceRequired = false;
        });

        // Edit Connection Modal
        document.getElementById('editConnectionClose').addEventListener('click', () => closeModal('editConnection'));
        document.getElementById('editConnectionBackdrop').addEventListener('click', () => closeModal('editConnection'));
        document.getElementById('saveConnectionBtn').addEventListener('click', () => {
            if(state.edgeBeingEdited) {
                state.edgeBeingEdited.text = document.getElementById('editChoiceText').value.trim() || "Choice";
                renderAll();
                autoSave();
            }
            closeModal('editConnection');
        });
        document.getElementById('deleteConnectionBtn').addEventListener('click', () => {
            if(state.edgeBeingEdited) {
                state.edges = state.edges.filter(e => e.id !== state.edgeBeingEdited.id);
                renderAll();
                autoSave();
            }
            closeModal('editConnection');
        });

        // Cancel connection mode on escape
        window.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && state.connectionMode) {
                cancelConnectionMode();
            }
        });
    }

    // --- Graph Logic ---

    function generateId() {
        return 'temp_' + Date.now() + '_' + Math.floor(Math.random() * 1000);
    }

    function createNode(name, type, x, y) {
        const node = {
            id: generateId(),
            name,
            type,
            x: Math.round(x / 20) * 20,
            y: Math.round(y / 20) * 20
        };
        state.nodes.push(node);
        renderAll();
        autoSave();
    }

    function deleteNode(id) {
        const node = state.nodes.find(n => n.id == id);
        if(!node || node.type === 'start') return; // Cannot delete start

        if(confirm(`Delete scene "${node.name}"?`)) {
            // Remove node
            state.nodes = state.nodes.filter(n => n.id != id);
            // Remove connected edges
            state.edges = state.edges.filter(e => e.from != id && e.to != id);
            
            renderAll();
            autoSave();
        }
    }

    function startConnectionMode(sourceId) {
        state.connectionMode = true;
        state.connectionSourceNodeId = sourceId;
        document.body.classList.add('connection-mode');
        
        // Create temporary drawing line
        state.drawingLine = document.createElementNS("http://www.w3.org/2000/svg", "path");
        state.drawingLine.setAttribute('class', 'drawing-line');
        renderEdges();
    }

    function completeConnection(targetId) {
        if (state.connectionSourceNodeId == targetId) {
            cancelConnectionMode();
            return;
        }

        const targetNode = state.nodes.find(n => n.id == targetId);
        if (targetNode.type === 'start') {
            alert("Cannot connect to Start node.");
            cancelConnectionMode();
            return;
        }

        // Check if already connected
        const exists = state.edges.some(e => e.from == state.connectionSourceNodeId && e.to == targetId);
        if (exists) {
            alert("Connection already exists.");
            cancelConnectionMode();
            return;
        }

        const sourceNode = state.nodes.find(n => n.id == state.connectionSourceNodeId);

        // Start node: no choice text — create edge directly, skip modal
        if (sourceNode && sourceNode.type === 'start') {
            createEdge(state.connectionSourceNodeId, targetId, '');
            return;
        }

        // Determine if choice text is required (source already has outgoing connections)
        const existingOut    = state.edges.filter(e => e.from == state.connectionSourceNodeId).length;
        const choiceRequired = existingOut >= 1;
        state._choiceRequired = choiceRequired;

        state.connectionTargetNodeId = targetId;
        const connectionInput = document.getElementById('connectionText');
        connectionInput.value = '';
        connectionInput.style.borderColor = '';
        connectionInput.placeholder = 'What choice leads here?';

        // Update label indicators
        const reqSpan = document.getElementById('connectChoiceRequired');
        const optSpan = document.getElementById('connectChoiceOptional');
        if (choiceRequired) {
            if (reqSpan) reqSpan.style.display = 'inline';
            if (optSpan) optSpan.style.display = 'none';
        } else {
            if (reqSpan) reqSpan.style.display = 'none';
            if (optSpan) optSpan.style.display = 'inline';
        }

        openModal('connectScene');
    }

    function cancelConnectionMode() {
        state.connectionMode = false;
        state.connectionSourceNodeId = null;
        state.connectionTargetNodeId = null;
        if (state.drawingLine) {
            state.drawingLine.remove();
            state.drawingLine = null;
        }
        document.body.classList.remove('connection-mode');
        renderEdges();
    }

    function createEdge(from, to, text) {
        state.edges.push({
            id: generateId(),
            from,
            to,
            text
        });
        cancelConnectionMode();
        renderAll();
        autoSave();
    }

    function openEditConnectionModal(edge) {
        state.edgeBeingEdited = edge;
        document.getElementById('editChoiceText').value = edge.text;
        openModal('editConnection');
    }

    // --- Modal Helpers ---
    function openModal(idPrefix) {
        document.getElementById(`${idPrefix}Modal`).classList.add('active');
        document.getElementById(`${idPrefix}Backdrop`).classList.add('active');
    }

    function closeModal(idPrefix) {
        document.getElementById(`${idPrefix}Modal`).classList.remove('active');
        document.getElementById(`${idPrefix}Backdrop`).classList.remove('active');
    }

    // Run
    init();
});
