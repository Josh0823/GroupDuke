<script lang="ts">
	import '../../static/main.css';

	import Form from '../components/Form.svelte';
	import FormInput from '../components/FormInput.svelte';
	import { isUserLoggedIn, serverURL } from '$lib/utils';
	import { onMount } from 'svelte';
	import TitleBar from '../components/TitleBar.svelte';

	let error: string;
	let id = 'add-course-form';
	let loggedIn: boolean;
	let netID: string;

	const addCourse = () => {
		let form = document.forms[id];
		const data = new FormData(form);

		let course = { id: 0 };
		for (let entry of data.entries()) {
			course[entry[0]] = entry[1];
		}
		course['user'] = netID;

		fetch(`${serverURL}/add-course`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify(course)
		}).then((res) => {
			if (res.ok) {
				window.location.replace('/');
			} else {
				error = 'Error adding course';
			}
		});
	};

	onMount(async () => {
		[loggedIn, netID] = isUserLoggedIn();

		if (loggedIn) {
		} else {
			window.location.replace('/login');
		}
	});
</script>

<main>
	<TitleBar {loggedIn} {netID} disabled={true} />

	<Form
		{id}
		title="Add a new course"
		{error}
		submitFn={addCourse}
		cancelFn={() => window.location.replace('/')}
	>
		<FormInput id="term" title="Term" placeholder="Sp22 | Fa21" />
		<FormInput id="course_number" title="Course Number" placeholder="BIO 201 | CHEM 101" />
		<FormInput id="professor" title="Professor" placeholder="John Smith | Jane Doe" />
		<FormInput id="time" title="Time" placeholder="MWF 8:00am | TTh 1:00pm" />
		<FormInput id="link" title="Link" placeholder="groupme.com/join/abc" />
	</Form>
</main>
