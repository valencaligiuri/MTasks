import {ref} from "vue";

const tasks = ref([]);

function fetchTasks(){
    fetch(`http://${window.location.hostname}:8081/api/tasks`, {
        method: "GET",
        headers: {'Content-Type': 'application/json'},
        credentials: "include"
    }).then(async res => {
        const data = await res.json();
        if (!res.ok) {
            console.error("Error")
        } else {
            tasks.value = data.tasks;
        }
    }).catch(() => console.error("error fetching tasks"))
}

export function useTasks() {
    return {
        tasks,
        fetchTasks
    }
}