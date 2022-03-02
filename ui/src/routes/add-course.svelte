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

	const checkTermFormat = (term: string) : boolean => {
		if (term.length != 4) {
			error = 'Term should only be four characters (i.e. Sp22)';
			return false;
		}

		let semester = term.substring(0, 2);
		if (!['Fa', 'Su', 'Sp'].includes(semester)) {
			error = 'First two digits of term should be Sp, Su, or Fa';
			return false;
		}

		let year = term.substring(2, 4);
		if (isNaN(parseFloat(year))) {
			error = 'Last two digits of term should be the year';
			return false;
		}

		return true;
	}

	const checkLinkFormat = (link: string) : boolean => {
		if (!link.startsWith('https://')) {
			error = 'Link should start with "https://"';
			return false;
		}

		if (!link.startsWith('https://groupme.com/join_group/')) {
			error = 'Link should start with "https://groupme.com/join_group/';
			return false;
		}

		return true;
	}

	const checkTimeFormat = (time: string) : boolean => {
		return true;
	}

	const addCourse = async () => {
		let form = document.forms[id];
		const data = new FormData(form);

		let course = { id: 0 };
		for (let entry of data.entries()) {
			course[entry[0]] = entry[1];
		}
		course['user'] = netID;

		if (!checkTermFormat(course['term'])) return;

		if (!checkLinkFormat(course['link'])) return;

		if (!checkTimeFormat(course['time'])) return;

		const res = await fetch(`${serverURL}/add-course`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify(course)
		});

		if (res.ok) {
			window.location.assign('/');
		} else {
			error = 'Error adding course';
		}
	};

	onMount(async () => {
		[loggedIn, netID] = isUserLoggedIn();

		if (loggedIn) {
		} else {
			window.location.replace('/login');
		}
	});

	let linkTitle = 'Link <small><a href="/help">(Where to find)</a></small>'
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
		<FormInput id="link" title={linkTitle} placeholder="groupme.com/join/abc" />
	</Form>
</main>
