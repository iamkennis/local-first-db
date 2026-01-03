export async function openDB() {
  return new Promise((resolve) => {
    const req = indexedDB.open("localdb", 1)
    req.onupgradeneeded = () =>
      req.result.createObjectStore("ops", { keyPath: "id" })
    req.onsuccess = () => resolve(req.result)
  })
}

export async function appendOp(op) {
  const db = await openDB()
  const tx = db.transaction("ops", "readwrite")
  tx.objectStore("ops").put(op)
}

export async function loadOps() {
  const db = await openDB()
  const tx = db.transaction("ops", "readonly")
  const store = tx.objectStore("ops")

  return new Promise(resolve => {
    const ops = []
    store.openCursor().onsuccess = e => {
      const cursor = e.target.result
      if (cursor) {
        ops.push(cursor.value)
        cursor.continue()
      } else {
        resolve(ops)
      }
    }
  })
}