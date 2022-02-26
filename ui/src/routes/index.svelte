<script lang="ts">
	// import css files
	import '../../static/main.css';
	import 'gridjs/dist/theme/mermaid.css';

	import { onMount } from 'svelte';
	import { writable } from 'svelte/store';

	// import npm components
	import { BaseComponent, h, html, PluginPosition } from 'gridjs';
	import Grid from 'gridjs-svelte';
	import Modal, { bind } from 'svelte-simple-modal';

	// import svelte components
	import LoginModal from '../components/LoginModal.svelte';
	import TitleBar from '../components/TitleBar.svelte';

	let serverURL = 'http://localhost:3000';
	let loggedIn = false;
	let netID = 'jmg136';
	let data = [];
	let showAllSemesters = false;
	let currentSemester = '';
	let grid: any;

	const columns = ['Term', 'Course Number', 'Professor', 'Time', 'Link'];

	const getCookie = (name: string) => {
		const value = `; ${document.cookie}`;
		const parts = value.split(`; ${name}=`);
		return parts.length === 2 ? parts.pop().split(';').shift() : '';
	};

	const deleteCookie = (name: string) => {
		document.cookie = name + '=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;';
	};

	const getCurrentSemester = () => {
		const d = new Date();
		const month = d.getMonth();
		const year = d.getFullYear();

		let term: string;
		if (month <= 4) {
			term = 'Sp';
		} else if (month <= 7) {
			term = 'Su';
		} else {
			term = 'Fa';
		}

		currentSemester = term.concat(`${year % 1000}`);
	};

	// Turns fetched data into valid format for datatable
	const formatData = (data: any[]) => {
		let ret = [];
		data.forEach((entry) => {
			entry.courseNumber = entry.course_number;
			entry.link = html(`<a href="${entry.link}" target="_blank">Join Group</a>`);

			if (showAllSemesters || entry.term == currentSemester) {
				ret.push(entry);
			}
		});

		return ret;
	};

	const getData = async () => {
		let url = `${serverURL}/data`;
		if (!showAllSemesters) {
			url += `?term=${currentSemester}`;
		}

		const response = await fetch(url, {
			method: 'GET',
			credentials: 'include',
			headers: {
				'Content-Type': 'text/plain'
			}
		});

		if (response.ok) {
			console.log('returning fetched and formatted data');
			data = formatData(await response.json());
			console.log(data);

			grid.updateConfig({ data: data }).forceRender();
		} else {
			console.error('Failed to fetch data');
			data = [];
		}
	};

	const handleLogin = (netID) => {
		loggedIn = true;
		netID = netID;

		getData();
	};

	const logout = () => {
		fetch(`${serverURL}/logout`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify(getCookie('session_token'))
		}).then((res) => {
			if (res.ok) {
				loggedIn = false;
				netID = '';

				deleteCookie('session_token');
				deleteCookie('net_id');
			} else {
				console.error('Error: logout failed');
			}
		});
	};

	const handleShowAllSemesters = () => {
		showAllSemesters = !showAllSemesters;
		getData();
	};

	onMount(async () => {
		// loggedIn = getCookie('session_token') != '';
		if (loggedIn) {
			getData();
		}
		getCurrentSemester();
	});

	const loginModal = writable(null);
	// @ts-ignore: Doesn't support type checking
	const showLoginModal = () => loginModal.set(bind(LoginModal, { serverURL: serverURL }, {onClose: () => console.log('closing')}));

	const registerModal = writable(null);
	// @ts-ignore: Doesn't support type checking
	const showRegisterModal = () => registerModal.set(bind(LoginModal));

	const addCourseModal = writable(null);
	// @ts-ignore: Doesn't support type checking
	const showAddCourseModal = () => addCourseModal.set(bind(LoginModal));

	class ButtonRowPlugin extends BaseComponent {
		render() {
			return h('div', { class: 'material-icons-row' }, [
				h('span', { class: 'material-icons', onclick: showAddCourseModal }, 'add'),
				h('span', { class: 'material-icons', onclick: null }, 'delete'),
				h(
					'span',
					{ class: 'material-icons', onclick: handleShowAllSemesters },
					`${showAllSemesters ? 'visibility' : 'visibility_off'}`
				)
			]);
		}
	}
</script>

<svelte:head>
	<title>GroupDuke</title>
</svelte:head>
<div>
	<Modal show={$loginModal} />
	<Modal show={$registerModal} />
	<Modal show={$addCourseModal} />

	<TitleBar
		{loggedIn}
		{netID}
		on:register={showRegisterModal}
		on:login={showLoginModal}
		on:logout={logout}
	/>

	<div class="content">
		{#if loggedIn}
			<Grid
				bind:instance={grid}
				{columns}
				{data}
				search={{ enabled: true }}
				pagination={{ enabled: true }}
				fixedHeader={true}
				sort={true}
				language={{ search: { placeholder: 'Search for a course' } }}
				plugins={[
					{
						id: 'button-row-plugin',
						component: ButtonRowPlugin,
						position: PluginPosition.Header
					}
				]}
			/>
		{:else}
			<p>Please login to view</p>
		{/if}
	</div>
</div>

<style>
	.content {
		align-content: center;
		text-align: center;
		margin: auto;
		padding: 20px 20px 20px 20px;
		max-width: 800px;
	}
</style>
