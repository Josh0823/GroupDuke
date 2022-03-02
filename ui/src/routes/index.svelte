<script lang="ts">
	import '../../static/main.css';
	import 'gridjs/dist/theme/mermaid.css';

	import { BaseComponent, h, html, PluginPosition } from 'gridjs';
	import {
		deleteCookie,
		getCookie,
		getCurrentSemester,
		isUserLoggedIn,
		serverURL
	} from '$lib/utils';
	import FooterBar from '../components/FooterBar.svelte';
	import Grid from 'gridjs-svelte';
	import { onMount } from 'svelte';
	import TitleBar from '../components/TitleBar.svelte';

	let grid: any;
	const columns = ['Term', 'Course Number', 'Professor', 'Time', 'Link'];
	let data = [];

	let loggedIn: boolean = true;
	let netID: string;
	let showAllSemesters = false;
	let currentSemester: string;

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

		const res = await fetch(url, {
			method: 'GET',
			credentials: 'include',
			headers: {
				'Content-Type': 'text/plain'
			}
		});

		if (res.ok) {
			data = formatData(await res.json());
			grid.updateConfig({ data: data }).forceRender();
		} else {
			console.error('Failed to fetch data');
			data = [];
		}
	};

	const handleLogout = async () => {
		const res = await fetch(`${serverURL}/logout`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify(getCookie('session_token'))
		});

		if (res.ok) {
			loggedIn = false;
			netID = '';

			deleteCookie('session_token');
			deleteCookie('net_id');

			window.location.assign('/login');
		} else {
			console.error('Error: logout failed');
		}
	};

	const handleShowAllSemesters = () => {
		showAllSemesters = !showAllSemesters;
		getData();
	};

	onMount(async () => {
		[loggedIn, netID] = isUserLoggedIn();

		if (loggedIn) {
			currentSemester = getCurrentSemester();
			getData();
		} else {
			window.location.replace('/login');
		}
	});

	class ButtonRowPlugin extends BaseComponent {
		render() {
			return h('div', { class: 'material-icons-row' }, [
				h(
					'span',
					{
						class: 'material-icons',
						title: 'Add a course',
						onclick: () => window.location.replace('/add-course')
					},
					'add'
				),
				h('span', { class: 'material-icons', title: 'Delete a course' }, 'delete'),
				h(
					'span',
					{ class: 'material-icons', title: 'Show all semesters', onclick: handleShowAllSemesters },
					`${showAllSemesters ? 'visibility' : 'visibility_off'}`
				)
			]);
		}
	}
</script>

<main>
	<TitleBar {loggedIn} {netID} on:logout={handleLogout} />

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
		{/if}
	</div>

	<FooterBar />
</main>

<style>
	/* @import "https://cdn.jsdelivr.net/npm/gridjs/dist/theme/mermaid.min.css"; */

	.content {
		align-content: center;
		text-align: center;
		margin: auto;
		padding: 20px 20px 20px 20px;
		max-width: 800px;
	}
</style>
