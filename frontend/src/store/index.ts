import { writable, get } from 'svelte/store';

type ChatEvent = {
    type: string;
    body: Record<string, any>;
};


function createStore() {
    const initialData = {
        status: null,
        isAuth: false,
        profile: {
            id: -1,
            sessionId: -1,
            username: null
        },
        events: []
    }

    const store = writable({ ...initialData });
    let ws: WebSocket;

    function onWSMessage(evt: MessageEvent) {
        const data = JSON.parse(evt.data)
        if (data.type === "login") {
            store.update(n => ({
                ...n,
                profile: {
                    id: data.body.userId,
                    sessionId: data.body.sessionId,
                    username: data.body.username,
                },
                isAuth: data.body.ok,
                events: [...n.events, data],
            }))
            return;
        }
        store.update(n => ({
            ...n,
            events: [
                ...n.events,
                data,
            ]
        }))

    }

    function onWSOpen() {
        store.update(n => ({ ...n, status: "ready" }))
    }

    function onWSClose() {
        store.update(n => ({ ...n, status: "close" }))
    }

    function send(data: ChatEvent) {
        ws.send(JSON.stringify(data))
    }

    return {
        subscribe: store.subscribe,
        connect() {
            ws = new WebSocket("ws://localhost:8080");

            ws.onopen = onWSOpen;
            ws.onmessage = onWSMessage;
            ws.onclose = onWSClose;
        },
        disconnect() {
            if (ws && ws.OPEN) {
                ws.close()
            }
        },
        send,
        login(username: string) {
            send({
                type: "login",
                body: { username, }
            })
        },
        getProfile(): Record<string, any> {
            return get(store).profile
        },
        reset: () => store.set({
            ...initialData
        })
    };
}

export const chat = createStore()

export const contextKey = Symbol("ctx")