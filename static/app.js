const form = document.querySelector('#new-alias');
const deleteForm = document.querySelector('#delete-alias');
const button = document.querySelector('#list-aliases');
const statusEl = document.querySelector('#status');
const aliases = document.querySelector('#aliases');

form.addEventListener('submit', async (event) => {
  event.preventDefault();
  form.querySelector('button').disabled = true;

  try {
    const response = await fetch('/new', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        alias: form.elements.alias.value,
        url: form.elements.url.value,
      }),
    });
    if (!response.ok) throw new Error('Unable to add alias.');

    statusEl.textContent = 'Alias added.';
    form.reset();
    button.click();
  } catch (error) {
    statusEl.textContent = error.message;
  } finally {
    form.querySelector('button').disabled = false;
  }
});

deleteForm.addEventListener('submit', async (event) => {
  event.preventDefault();
  deleteForm.querySelector('button').disabled = true;

  try {
    const response = await fetch(`/${encodeURIComponent(deleteForm.elements.alias.value)}`, {
      method: 'DELETE',
    });
    if (!response.ok) throw new Error('Unable to delete alias.');

    statusEl.textContent = 'Alias deleted.';
    deleteForm.reset();
    button.click();
  } catch (error) {
    statusEl.textContent = error.message;
  } finally {
    deleteForm.querySelector('button').disabled = false;
  }
});

button.addEventListener('click', async () => {
  button.disabled = true;
  statusEl.textContent = 'Loading...';
  aliases.replaceChildren();

  try {
    const response = await fetch('/list');
    if (!response.ok) throw new Error('Unable to load aliases.');

    const entries = Object.entries(await response.json()).sort();
    for (const [alias] of entries) {
      const link = document.createElement('a');
      link.href = `/${encodeURIComponent(alias)}`;
      link.textContent = alias;

      const item = document.createElement('li');
      item.append(link);
      aliases.append(item);
    }
    statusEl.textContent = `${entries.length} aliases`;
  } catch (error) {
    statusEl.textContent = error.message;
  } finally {
    button.disabled = false;
  }
});
